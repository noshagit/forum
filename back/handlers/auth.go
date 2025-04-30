package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Pseudo          string `json:"pseudo"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func LoginHandler(router *mux.Router) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		http.ServeFile(w, r, "") // TODO : File path
	}).Methods("GET")

	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		http.ServeFile(w, r, "") // TODO : File path
	}).Methods("GET")

	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		http.ServeFile(w, r, "") // TODO : File path
	}).Methods("GET")

	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "JSON decoding error", http.StatusBadRequest)
			log.Println("JSON decoding error:", err)
			return
		}

		db, err := sql.Open("sqlite3", "") // TODO : DB file path
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			log.Println("Database connection error:", err)
			return
		}
		defer db.Close()

		var storedHashedPassword string
		row := db.QueryRow("SELECT password FROM users WHERE email = ?", credentials.Email)
		err = row.Scan(&storedHashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database query error", http.StatusInternalServerError)
				log.Println("Database query error:", err)
			}
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(credentials.Password)) == nil {
			sessionToken := uuid.New().String()

			_, err := db.Exec("INSERT INTO sessions (token, email) VALUES (?, ?)", sessionToken, credentials.Email)
			if err != nil {
				http.Error(w, "Error creating session", http.StatusInternalServerError)
				log.Println("Session insertion error:", err)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:   "session_token",
				Value:  sessionToken,
				Path:   "/",
				MaxAge: 86400,
			})

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Login successful"))
		} else {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		}
	}).Methods("POST")
}

func ProfileHandler(router *mux.Router) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		http.ServeFile(w, r, "") // TODO : File path
	}).Methods("GET")

	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		http.ServeFile(w, r, "") // TODO : File path
	}).Methods("GET")

	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		http.ServeFile(w, r, "") // TODO : File path
	}).Methods("GET")

	router.HandleFunc("/api/get-profile", func(w http.ResponseWriter, r *http.Request) { 
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		db, err := sql.Open("sqlite3", "") // TODO : DB file path
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var email string
		row := db.QueryRow("SELECT email FROM sessions WHERE token = ?", cookie.Value)
		err = row.Scan(&email)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		var profile struct {
			Pseudo string `json:"pseudo"`
			Email  string `json:"email"`
		}

		row = db.QueryRow("SELECT pseudo, email FROM users WHERE email = ?", email)
		err = row.Scan(&profile.Pseudo, &profile.Email)
		if err != nil {
			http.Error(w, "Error retrieving profile information", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"profile": profile,
		})
	}).Methods("GET")

}

func RegisterHandler(router *mux.Router) {
	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		http.ServeFile(w, r, "") // TODO : File path
	}).Methods("GET")

	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		http.ServeFile(w, r, "") // TODO : File path
	}).Methods("GET")

	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		http.ServeFile(w, r, "") // TODO : File path
	}).Methods("GET")

	router.HandleFunc("", func(w http.ResponseWriter, r *http.Request) { // TODO : Route path
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "JSON decoding error", http.StatusBadRequest)
			log.Println("JSON decoding error:", err)
			return
		}

		if user.Pseudo == "" || user.Email == "" || user.Password == "" || user.ConfirmPassword == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if user.Password != user.ConfirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		hashedPassword := hashPassword(user.Password)

		db, err := sql.Open("sqlite3", "") // TODO : DB file path
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			log.Println("Database connection error:", err)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO users (pseudo, email, password) VALUES (?, ?, ?)")
		if err != nil {
			log.Println("Database preparation error:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		if _, err = stmt.Exec(user.Pseudo, user.Email, hashedPassword); err != nil {
			log.Println("Database insertion error:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Registration successful"))
	}).Methods("POST")
}

func LogoutHandler(router *mux.Router) {
	router.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err == nil {
			db, err := sql.Open("sqlite3", "") // TODO : DB file path
			if err == nil {
				_, err := db.Exec("DELETE FROM sessions WHERE token = ?", cookie.Value)
				if err != nil {
					http.Error(w, "Error deleting session", http.StatusInternalServerError)
					log.Println("Error deleting session:", err)
				}
				db.Close()
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		http.Redirect(w, r, "/main-menu/menu.html", http.StatusSeeOther)
	}).Methods("POST")
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func CheckPassword(password, hashedPassword string) bool {
	hashedInput := hashPassword(password)
	return hashedInput == hashedPassword
}
