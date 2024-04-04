package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    _ "github.com/mattn/go-sqlite3"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

var db *sql.DB

func main() {
    // Initialize the database
    var err error
    db, err = sql.Open("sqlite3", "./test.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create user table
    createUserTable()

    // Create a new router instance
    router := mux.NewRouter()

    // Middleware for logging
    router.Use(loggingMiddleware)

    // Define API routes
    router.HandleFunc("/users", getUsers).Methods("GET")
    router.HandleFunc("/users/{id}", getUser).Methods("GET")
    router.HandleFunc("/users", createUser).Methods("POST")
    router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
    router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

    // Start the server
    fmt.Println("Server is running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}

func createUserTable() {
    query := `
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT,
            email TEXT
        );
    `
    _, err := db.Exec(query)
    if err != nil {
        log.Fatal(err)
    }
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    // Fetch all users from the database
    rows, err := db.Query("SELECT id, username, email FROM users")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
            log.Fatal(err)
        }
        users = append(users, user)
    }

    // Convert users slice to JSON
    json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    // Get user ID from request params
    vars := mux.Vars(r)
    id := vars["id"]

    // Fetch user from the database
    var user User
    row := db.QueryRow("SELECT id, username, email FROM users WHERE id = ?", id)
    if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
        log.Fatal(err)
    }

    // Convert user struct to JSON
    json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User

    // Decode JSON request body into User struct
    json.NewDecoder(r.Body).Decode(&user)

    // Insert user into the database
    result, err := db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", user.Username, user.Email)
    if err != nil {
        log.Fatal(err)
    }

    // Get the ID of the inserted user
    id, _ := result.LastInsertId()

    // Return the ID of the newly created user
    fmt.Fprintf(w, "New user created with ID: %d", id)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    // Get user ID from request params
    vars := mux.Vars(r)
    id := vars["id"]

    // Decode JSON request body into User struct
    var user User
    json.NewDecoder(r.Body).Decode(&user)

    // Update user in the database
    _, err := db.Exec("UPDATE users SET username = ?, email = ? WHERE id = ?", user.Username, user.Email, id)
    if err != nil {
        log.Fatal(err)
    }

    // Return success message
    fmt.Fprintf(w, "User with ID %s has been updated", id)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    // Get user ID from request params
    vars := mux.Vars(r)
    id := vars["id"]

    // Delete user from the database
    _, err := db.Exec("DELETE FROM users WHERE id = ?", id)
    if err != nil {
        log.Fatal(err)
    }

    // Return success message
    fmt.Fprintf(w, "User with ID %s has been deleted", id)
}
