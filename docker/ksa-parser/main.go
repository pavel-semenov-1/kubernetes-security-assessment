package main

import (
	"encoding/json"
	"fmt"
	"ksa-parser/parser"
	"net/http"
	"os"
)

func main() {
	// Initialize parsers
	trivyParser := parser.NewTrivyParser()
	kubeBenchParser := parser.NewKubeBenchParser()

	// Parse data once during startup
	trivyDataLocation := os.Getenv("TRIVY_DATA_LOCATION")
	if trivyDataLocation == "" {
		panic("TRIVY_DATA_LOCATION environment variable not set")
	}
	if err := trivyParser.Parse(trivyDataLocation); err != nil {
		fmt.Println("Error parsing Trivy data:", err)
		return
	}
	kubeBenchDataLocation := os.Getenv("KUBE_BENCH_DATA_LOCATION")
	if kubeBenchDataLocation == "" {
		panic("KUBE_BENCH_DATA_LOCATION environment variable not set")
	}
	if err := kubeBenchParser.Parse("/Users/pavel/Diploma/kubernetes-security-assessment/artifacts/kube-bench-scan.json"); err != nil {
		fmt.Println("Error parsing Kube-bench data:", err)
		return
	}

	parsers := map[string]parser.Parser{
		"Trivy":      trivyParser,
		"Kube-bench": kubeBenchParser,
	}

	// HTTP handler to query results
	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		if scanner == "" || parsers[scanner] == nil {
			http.Error(w, "Invalid or missing scanner parameter", http.StatusBadRequest)
			return
		}

		results := parsers[scanner].GetResults()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	// HTTP handler to query vulnerabilities by severity
	http.HandleFunc("/vulnerabilities", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		namespace := r.URL.Query().Get("namespace")
		severity := r.URL.Query().Get("severity")

		if scanner == "" || parsers[scanner] == nil {
			http.Error(w, "Missing scanner parameter", http.StatusBadRequest)
			return
		}

		results := parsers[scanner].GetVulnerabilities(&namespace, &severity)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	// HTTP handler to query misconfigurations by severity
	http.HandleFunc("/misconfigurations", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		namespace := r.URL.Query().Get("namespace")
		severity := r.URL.Query().Get("severity")

		if scanner == "" || parsers[scanner] == nil {
			http.Error(w, "Missing scanner parameter", http.StatusBadRequest)
			return
		}

		results := parsers[scanner].GetMisconfigurations(&namespace, &severity)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
