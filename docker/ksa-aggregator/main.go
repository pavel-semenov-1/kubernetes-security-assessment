package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"ksa-aggregator/runner"
	"net/http"
	"os"
)

func main() {
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		panic("NAMESPACE environment variable not set")
	}

	// Create Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		panic("Failed to load in-cluster config: " + err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("Failed to create Kubernetes client: " + err.Error())
	}

	trivyRunner := runner.NewTrivyRunner(clientset, namespace)
	runners := map[string]runner.Runner{
		"Trivy": trivyRunner,
	}

	// HTTP handler to run scans
	http.HandleFunc("/scan", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		if scanner == "" || runners[scanner] == nil {
			http.Error(w, "Invalid or missing scanner parameter", http.StatusBadRequest)
			return
		}

		jobStatus := runners[scanner].GetStatus()
		if jobStatus.Active() {
			http.Error(w, "Job is already running", http.StatusBadRequest)
			return
		}

		runners[scanner].Run()
	})

	// HTTP handler to query scan status
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		if scanner == "" || runners[scanner] == nil {
			http.Error(w, "Invalid or missing scanner parameter", http.StatusBadRequest)
			return
		}

		jobStatus := runners[scanner].GetStatus()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(runner.GetJobState(jobStatus))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server running on port", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
