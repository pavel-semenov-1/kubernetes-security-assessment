package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"ksa-parser/parser"
	"ksa-parser/pdb"
	"net/http"
	"os"
	"sync"
)

func main() {
	// Get environment variables
	reportDataLocation := os.Getenv("REPORT_DATA_LOCATION")
	if reportDataLocation == "" {
		panic("REPORT_DATA_LOCATION environment variable not set")
	}
	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		panic("POSTGRES_HOST environment variable not set")
	}
	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		panic("POSTGRES_PORT environment variable not set")
	}
	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		panic("POSTGRES_USER environment variable not set")
	}
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		panic("POSTGRES_PASSWORD environment variable not set")
	}
	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		panic("POSTGRES_DB environment variable not set")
	}

	// Initialize the database connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	dbcon, err := pdb.NewDB(connStr)
	if err != nil {
		panic(err)
	}
	defer dbcon.Close()

	// Initialize parsers
	trivyParser := parser.NewTrivyParser()
	kubeBenchParser := parser.NewKubeBenchParser()
	prowlerParser := parser.NewProwlerParser()
	kubescapeParser := parser.NewKubescapeParser()

	parsers := map[string]parser.Parser{
		"trivy":      trivyParser,
		"kube-bench": kubeBenchParser,
		"prowler":    prowlerParser,
		"kubescape":  kubescapeParser,
	}

	parseAndPopulate := func(reportId string, w *http.ResponseWriter) error {
		fmt.Println("Getting scanner name...")
		scanner, err := dbcon.GetScannerNameByReportId(reportId)
		if err != nil {
			return fmt.Errorf("error getting the scanner name %v", err)
		}
		fmt.Println("Querying the report filename...")
		filename := dbcon.GetReportFilename(reportId)
		filepath := fmt.Sprintf("%s/%s/%s", reportDataLocation, scanner, filename)
		fmt.Printf("Parsing the report file %s...\n", filepath)
		vuln, misc, err := parsers[scanner].Parse(filepath)
		if err != nil {
			return fmt.Errorf("error parsing the report: %v", err)
		}
		fmt.Println("Populating the database...")
		err = dbcon.Populate(reportId, vuln, misc)
		if err != nil {
			return fmt.Errorf("error populating the database: %v", err)
		}
		return nil
	}

	queryDBMiscVulnData := func(reportId string, search string, resolved bool, w *http.ResponseWriter) ([]parser.Misconfiguration, []parser.Vulnerability, error) {
		var reportIds []string

		if reportId == "" {
			reportIds, err = dbcon.GetLastParsedReportIds()
			if err != nil {
				return nil, nil, fmt.Errorf("error getting the last parsed report id %v", err)
			}
			if len(reportIds) == 0 {
				return []parser.Misconfiguration{}, []parser.Vulnerability{}, nil
			}
		} else {
			reportIds = []string{reportId}
		}

		misconfigurations, err := dbcon.GetMisconfigurations(reportIds, search, resolved)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting misconfigurations %v", err)
		}
		vulnerabilities, err := dbcon.GetVulnerabilities(reportIds, search)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting vulnerabilities %v", err)
		}

		return misconfigurations, vulnerabilities, nil
	}

	// Websocket stuff
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	var clients = make(map[*websocket.Conn]bool)
	var mutex = sync.Mutex{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("WebSocket upgrade failed:", err)
			return
		}
		defer conn.Close()

		// Register the client
		mutex.Lock()
		clients[conn] = true
		mutex.Unlock()

		// Listen for disconnects
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				mutex.Lock()
				delete(clients, conn)
				mutex.Unlock()
				break
			}
		}
	})

	broadcastMessage := func(message string) {
		mutex.Lock()
		defer mutex.Unlock()

		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("Failed to send message to client:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}

	http.HandleFunc("/broadcast", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.Message == "" {
			http.Error(w, "Empty message", http.StatusBadRequest)
			return
		}

		broadcastMessage(req.Message)
	})

	http.HandleFunc("/parse", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ReportId string `json:"report_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if req.ReportId == "" {
			http.Error(w, "Missing reportId", http.StatusBadRequest)
			return
		}

		broadcastMessage(fmt.Sprintf("@parse:%s:start", req.ReportId))

		err := parseAndPopulate(req.ReportId, &w)
		if err != nil {
			broadcastMessage(fmt.Sprintf("Error parsing the report: %s", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
		broadcastMessage(fmt.Sprintf("Report with id %s has been successfully parsed", req.ReportId))
		broadcastMessage(fmt.Sprintf("@parse:%s:end", req.ReportId))
	})

	// HTTP handler to query scanners
	http.HandleFunc("/scanners", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		scanners, err := dbcon.GetScannerNames()
		if err != nil {
			http.Error(w, "Error getting scanners", http.StatusInternalServerError)
			panic(err)
		}
		json.NewEncoder(w).Encode(scanners)
	})

	// HTTP handler to query vulnerabilities
	http.HandleFunc("/vulnerabilities", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			reportId := r.URL.Query().Get("reportId")
			search := r.URL.Query().Get("search")

			fmt.Printf("Querying vulnerabilities for reportId=%s and search=%s\n", reportId, search)
			_, vulnerabilities, err := queryDBMiscVulnData(reportId, search, false, &w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				panic(err)
			}

			vulnMap := make(map[string][]parser.Vulnerability)
			for _, vuln := range vulnerabilities {
				vulnMap[vuln.Severity] = append(vulnMap[vuln.Severity], vuln)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(vulnMap)
		case http.MethodDelete:
			reportId := r.URL.Query().Get("reportId")
			id := r.URL.Query().Get("id")
			search := r.URL.Query().Get("search")

			fmt.Printf("Deleting vulnerabilities for reportId=%s, id=%s and search=%s\n", reportId, id, search)

			if id == "" {
				var reportIds []string
				if reportId == "" {
					reportIds, err = dbcon.GetLastParsedReportIds()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						panic(err)
					}
				} else {
					reportIds = []string{reportId}
				}
				err = dbcon.DeleteVulnerabilitiesBySearchTerm(reportIds, search)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					panic(err)
				}
			} else {
				err = dbcon.DeleteVulnerabilityById(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					panic(err)
				}
			}
		}
	})

	// HTTP handler to query misconfigurations
	http.HandleFunc("/misconfigurations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			reportId := r.URL.Query().Get("reportId")
			search := r.URL.Query().Get("search")
			resolved := r.URL.Query().Get("resolved")

			fmt.Printf("Querying misconfigurations for reportId=%s, search=%s and resolved=%s\n", reportId, search, resolved)

			misconfigurations, _, err := queryDBMiscVulnData(reportId, search, resolved == "true", &w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				panic(err)
			}

			miscMap := make(map[string][]parser.Misconfiguration)
			for _, misc := range misconfigurations {
				miscMap[misc.Severity] = append(miscMap[misc.Severity], misc)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(miscMap)
		case http.MethodDelete:
			reportId := r.URL.Query().Get("reportId")
			id := r.URL.Query().Get("id")
			search := r.URL.Query().Get("search")

			fmt.Printf("Deleting misconfigurations for reportId=%s, id=%s and search=%s\n", reportId, id, search)

			if id == "" {
				var reportIds []string
				if reportId == "" {
					reportIds, err = dbcon.GetLastParsedReportIds()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						panic(err)
					}
				} else {
					reportIds = []string{reportId}
				}
				err = dbcon.DeleteMisconfigurationsBySearchTerm(reportIds, search)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					panic(err)
				}
			} else {
				err = dbcon.DeleteMisconfigurationById(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					panic(err)
				}
			}
		case http.MethodPatch:
			reportId := r.URL.Query().Get("reportId")
			id := r.URL.Query().Get("id")
			search := r.URL.Query().Get("search")

			fmt.Printf("Resolving misconfigurations for reportId=%s, id=%s and search=%s\n", reportId, id, search)

			if id == "" {
				var reportIds []string
				if reportId == "" {
					reportIds, err = dbcon.GetLastParsedReportIds()
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						panic(err)
					}
				} else {
					reportIds = []string{reportId}
				}
				err = dbcon.ResolveMisconfigurationsBySearchTerm(reportIds, search)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					panic(err)
				}
			} else {
				err = dbcon.ResolveMisconfigurationById(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					panic(err)
				}
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// HTTP handler to query scan reports
	http.HandleFunc("/reports", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")

		if scanner == "" || parsers[scanner] == nil {
			http.Error(w, "Missing scanner parameter", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		reports, err := dbcon.GetReports(scanner)
		if err != nil {
			http.Error(w, "Error getting reports", http.StatusInternalServerError)
			panic(err)
		}
		json.NewEncoder(w).Encode(reports)
	})

	// HTTP handler to query statistics
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")

		if scanner == "" || parsers[scanner] == nil {
			http.Error(w, "Missing scanner parameter", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		reportId, err := dbcon.GetLastParsedReportId(scanner)
		if err != nil {
			http.Error(w, "Error querying reports", http.StatusInternalServerError)
			panic(err)
		}

		type Stats struct {
			Vulnerabilities   map[string]int `json:"vulnerabilities"`
			Misconfigurations map[string]int `json:"misconfigurations"`
		}

		vuln, err := dbcon.GetVulnerabilityStatistics(reportId)
		if err != nil {
			http.Error(w, "Error getting vulnerability statistics", http.StatusInternalServerError)
			panic(err)
		}

		misc, err := dbcon.GetMisconfigurationStatistics(reportId)
		if err != nil {
			http.Error(w, "Error getting misconfiguration statistics", http.StatusInternalServerError)
			panic(err)
		}

		result := Stats{
			Vulnerabilities:   vuln,
			Misconfigurations: misc,
		}

		json.NewEncoder(w).Encode(result)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(fmt.Errorf("error starting server: %v", err))
	}
}
