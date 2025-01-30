package main

import (
	"encoding/json"
	"fmt"
	"ksa-parser/parser"
	"ksa-parser/pdb"
	"net/http"
	"os"
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

	parsers := map[string]parser.Parser{
		"trivy":     trivyParser,
		"kubebench": kubeBenchParser,
	}

	parseAndPopulate := func(reportId string, scanner string, prsr *parser.Parser, w *http.ResponseWriter) error {
		fmt.Println("Querying the report filename...")
		filename := dbcon.GetReportFilename(reportId)
		filepath := fmt.Sprintf("%s/%s/%s", reportDataLocation, scanner, filename)
		fmt.Printf("Parsing the report file %s...\n", filepath)
		vuln, misc, err := (*prsr).Parse(filepath)
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

	// HTTP handler to query vulnerabilities by severity
	http.HandleFunc("/vulnerabilities", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		// namespace := r.URL.Query().Get("namespace")
		severity := r.URL.Query().Get("severity")
		reportId := r.URL.Query().Get("reportId")

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
			return
		}
		if !parsed {
			fmt.Println("Report not parsed, parsing...")
			err = parseAndPopulate(reportId, scanner, &prsr, &w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		fmt.Println("Report parsed.")

		w.Header().Set("Content-Type", "application/json")
		fmt.Println("Querying vulnerabilities from the database...")
		vulnerabilities, err := dbcon.GetVulnerabilities(reportId, severity)
		if err != nil {
			http.Error(w, "Error getting vulnerabilities from the database", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(vulnerabilities)
	})

	// HTTP handler to query misconfigurations by severity
	http.HandleFunc("/misconfigurations", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		// namespace := r.URL.Query().Get("namespace")
		severity := r.URL.Query().Get("severity")
		reportId := r.URL.Query().Get("reportId")

		prsr := parsers[scanner]

		if scanner == "" || prsr == nil {
			http.Error(w, "Missing scanner parameter", http.StatusBadRequest)
			return
		}

		parsed, err := dbcon.Parsed(reportId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !parsed {
			err = parseAndPopulate(reportId, scanner, &prsr, &w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		misconfigurations, err := dbcon.GetMisconfigurations(reportId, severity)
		if err != nil {
			http.Error(w, "Error getting misconfigurations from the database", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(misconfigurations)
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
			return
		}
		json.NewEncoder(w).Encode(reports)
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
