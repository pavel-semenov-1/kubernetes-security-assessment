package main

import (
	"encoding/json"
	"fmt"
	"ksa-parser/parser"
	"ksa-parser/pdb"
	"net/http"
	"os"
	"strconv"
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

	parsers := map[string]parser.Parser{
		"trivy":      trivyParser,
		"kube-bench": kubeBenchParser,
		"prowler":    prowlerParser,
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

		err := parseAndPopulate(req.ReportId, &w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
	})

	// HTTP handler to query vulnerabilities
	http.HandleFunc("/vulnerabilities", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		reportId := r.URL.Query().Get("reportId")
		search := r.URL.Query().Get("search")

		fmt.Printf("Querying vulnerabilities for reportId=%s and search=%s\n", reportId, search)

		prsr := parsers[scanner]

		if scanner == "" || prsr == nil {
			http.Error(w, "Missing scanner parameter", http.StatusBadRequest)
			return
		}
		if reportId == "" {
			http.Error(w, "Missing reportId parameter", http.StatusBadRequest)
			return
		}

		fmt.Println("Checking if this report already parsed...")
		parsed, err := dbcon.Parsed(reportId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			panic(err)
		}
		if !parsed {
			fmt.Println("Report not parsed, parsing...")
			err = parseAndPopulate(reportId, &w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				panic(err)
			}
		}
		fmt.Println("Report parsed.")

		w.Header().Set("Content-Type", "application/json")
		fmt.Println("Querying vulnerabilities from the database...")
		vulnerabilities, err := dbcon.GetVulnerabilities(reportId, search)
		if err != nil {
			http.Error(w, "Error getting vulnerabilities from the database", http.StatusInternalServerError)
			panic(err)
		}

		vulnMap := make(map[string][]parser.Vulnerability)
		for _, vuln := range vulnerabilities {
			vulnMap[vuln.Severity] = append(vulnMap[vuln.Severity], vuln)
		}

		json.NewEncoder(w).Encode(vulnMap)
	})

	// HTTP handler to query misconfigurations
	http.HandleFunc("/misconfigurations", func(w http.ResponseWriter, r *http.Request) {
		reportId := r.URL.Query().Get("reportId")
		search := r.URL.Query().Get("search")

		fmt.Printf("Querying misconfigurations for reportId=%s and search=%s\n", reportId, search)

		var reports []string

		if reportId == "" {
			for scn, _ := range parsers {
				rid, err := dbcon.GetLastParsedReportId(scn)
				if err != nil {
					http.Error(w, "Error getting reports from the database", http.StatusInternalServerError)
					panic(err)
				}
				reports = append(reports, strconv.Itoa(rid))
			}
		} else {
			reports = append(reports, reportId)
		}

		misconfigurations, err := dbcon.GetMisconfigurations(reports, search)
		if err != nil {
			http.Error(w, "Error getting misconfigurations from the database", http.StatusInternalServerError)
			panic(err)
		}

		miscMap := make(map[string][]parser.Misconfiguration)
		for _, misc := range misconfigurations {
			miscMap[misc.Severity] = append(miscMap[misc.Severity], misc)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(miscMap)
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
