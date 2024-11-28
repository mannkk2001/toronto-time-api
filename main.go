package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Response struct {
	CurrentTime string `json:"current_time"`
}

var db *sql.DB

func connectDatabase() {
	var err error
	// Update with your MySQL username, password, and database name
	dsn := "root:password23#@tcp(127.0.0.1:3306)/toronto_time_db"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}
	log.Println("Database connection established.")
}

func getCurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	// Get current time in Toronto
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Error loading timezone", http.StatusInternalServerError)
		return
	}
	currentTime := time.Now().In(loc)

	// Log the current time to the database
	_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", currentTime)
	if err != nil {
		http.Error(w, "Error logging time to database", http.StatusInternalServerError)
		return
	}

	// Respond with the current time in JSON format
	response := Response{CurrentTime: currentTime.Format("2006-01-02 15:04:05")}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Connect to the database
	connectDatabase()
	defer db.Close()

	// Set up the router and endpoints
	router := mux.NewRouter()
	router.HandleFunc("/current-time", getCurrentTimeHandler).Methods("GET")

	// Start the server
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
