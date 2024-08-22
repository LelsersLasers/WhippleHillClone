package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"sync"
	"time"
	// "encoding/json"
	_ "github.com/mattn/go-sqlite3"
)

const Port = 8080
const DbPath = "./database.db"

const SessionIdCookieName = "session_id"
const SessionTimeout = 24 * time.Hour

const ContextFailCookieNameBase = "context_fail_"
const ContextFailCookieTimeout = 5 * time.Second

var (
	db    *sql.DB
	mutex sync.Mutex
)

func main() {
	db = dbConn()
	createTables(db)
	defer db.Close()

	handler := http.NewServeMux()

	handler.HandleFunc("/", homePage)
	handler.HandleFunc("/login", loginPage)
	handler.HandleFunc("/register", registerPage)

	handler.HandleFunc("/login_user", loginUser)
	handler.HandleFunc("/logout_user", logoutUser)
	handler.HandleFunc("/register_user", registerUser)

	loggedHandler := loggingMiddleware(handler)

	fmt.Printf("Server is running on port %d\n", Port)

	addr := fmt.Sprintf(":%d", Port)
	err := http.ListenAndServe(addr, loggedHandler)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}
}
