package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
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
	trivyJobName := os.Getenv("TRIVY_JOB_NAME")
	if trivyJobName == "" {
		trivyJobName = "trivy-runner"
	}

	// Open a database connection
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	con, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Errorf("unable to open DB: %v", err))
	}
	defer con.Close()

	// Create Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		panic("Failed to load in-cluster config: " + err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("Failed to create Kubernetes client: " + err.Error())
	}

	trivyRunner := runner.NewTrivyRunner(clientset, namespace, trivyJobName)
	runners := map[string]runner.Runner{
		"trivy": trivyRunner,
	}

	// HTTP handler to run scans
	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		rnr := runners[scanner]
		if scanner == "" || rnr == nil {
			http.Error(w, "Invalid or missing scanner parameter", http.StatusBadRequest)
			return
		}

		jobStatus := rnr.GetStatus()
		if jobStatus.Active() {
			http.Error(w, "Job is already running", http.StatusBadRequest)
			return
		} else if jobStatus.Failed() || jobStatus.Succeeded() {
			err := rnr.CleanUp()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err := rnr.Run()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go rnr.Watch(con)
	})

	// HTTP handler to query scan status
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		scanner := r.URL.Query().Get("scanner")
		rnr := runners[scanner]
		if scanner == "" || rnr == nil {
			http.Error(w, "Invalid or missing scanner parameter", http.StatusBadRequest)
			return
		}

		jobStatus := rnr.GetStatus()
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(jobStatus)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
