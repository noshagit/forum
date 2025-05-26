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

type ChangePasswordRequest struct {
	Password string `json:"password"`
}

func LoginHandler(router *mux.Router) {
	router.HandleFunc("/front/login/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/login/login.html")
	}).Methods("GET")

	router.HandleFunc("/front/login/login.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/login/login.css")
	}).Methods("GET")

	router.HandleFunc("/front/login/login.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/login/login.js")
	}).Methods("GET")

	router.HandleFunc("/front/images/login.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/login.png")
	}).Methods("GET")

	router.HandleFunc("/front/login/login.html", func(w http.ResponseWriter, r *http.Request) {
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

		db, err := sql.Open("sqlite3", "../database/bddforum.db")
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

func RegisterHandler(router *mux.Router) {
	router.HandleFunc("/front/register/register.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/register/register.html")
	}).Methods("GET")

	router.HandleFunc("/front/register/register.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/register/register.css")
	}).Methods("GET")

	router.HandleFunc("/front/register/register.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/register/register.js")
	}).Methods("GET")

	router.HandleFunc("/front/images/register.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/register.png")
	}).Methods("GET")

	router.HandleFunc("/front/register/register.html", func(w http.ResponseWriter, r *http.Request) {
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
			log.Printf("Pseudo: %s, Email: %s, Password: %s, ConfirmPassword: %s", user.Pseudo, user.Email, user.Password, user.ConfirmPassword)
			return
		}

		if user.Password != user.ConfirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		hashedPassword := hashPassword(user.Password)

		db, err := sql.Open("sqlite3", "../database/bddforum.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			log.Println("Database connection error:", err)
			return
		}
		defer db.Close()

		stmt, err := db.Prepare("INSERT INTO users (username, email, password) VALUES (?, ?, ?)")
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

		sessionToken := uuid.New().String()

		_, err = db.Exec("INSERT INTO sessions (token, email) VALUES (?, ?)", sessionToken, user.Email)
		if err != nil {
			log.Println("Session insertion error:", err)
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  sessionToken,
			Path:   "/",
			MaxAge: 86400,
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Registration successful"))
		log.Println("POST reçu")
	}).Methods("POST")
}

func ProfileHandler(router *mux.Router) {
	router.HandleFunc("/front/profil/profil.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/profil/profil.html")
	}).Methods("GET")

	router.HandleFunc("/front/profil/profil.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/profil/profil.css")
	}).Methods("GET")

	router.HandleFunc("/front/profil/profil.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/profil/profil.js")
	}).Methods("GET")

	router.HandleFunc("/front/images/pfp.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/pfp.png")
	}).Methods("GET")

	router.HandleFunc("/front/images/tnt.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/tnt.png")
	}).Methods("GET")

	router.HandleFunc("/front/images/background.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/background.png")
	}).Methods("GET")

	router.HandleFunc("/get-profile", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		db, err := sql.Open("sqlite3", "../database/bddforum.db")
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
			Pseudo string `json:"username"`
			Email  string `json:"email"`
			ID     int    `json:"id"`
		}

		row = db.QueryRow("SELECT id, username, email FROM users WHERE email = ?", email)
		err = row.Scan(&profile.ID, &profile.Pseudo, &profile.Email)
		if err != nil {
			http.Error(w, "Error retrieving profile information", http.StatusInternalServerError)
			return
		}

		log.Println("Email from session:", email)
		log.Println("Username from session:", profile.Pseudo)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"profile": profile,
		})
		log.Println("POST reçu")
	}).Methods("GET")

	router.HandleFunc("/update-profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		var updateData struct {
			Username string `json:"username"`
			Email    string `json:"email"`
		}

		err = json.NewDecoder(r.Body).Decode(&updateData)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if updateData.Username == "" || updateData.Email == "" {
			http.Error(w, "Username and Email cannot be empty", http.StatusBadRequest)
			return
		}

		db, err := sql.Open("sqlite3", "../database/bddforum.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var currentEmail string
		err = db.QueryRow("SELECT email FROM sessions WHERE token = ?", cookie.Value).Scan(&currentEmail)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		_, err = db.Exec("UPDATE users SET username = ?, email = ? WHERE email = ?", updateData.Username, updateData.Email, currentEmail)
		if err != nil {
			http.Error(w, "Error updating profile", http.StatusInternalServerError)
			return
		}

		if updateData.Email != currentEmail {
			_, err = db.Exec("UPDATE sessions SET email = ? WHERE token = ?", updateData.Email, cookie.Value)
			if err != nil {
				http.Error(w, "Error updating session email", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Profile updated successfully"))
	}).Methods("POST")

	router.HandleFunc("/delete-profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Not authenticated", http.StatusUnauthorized)
			return
		}

		db, err := sql.Open("sqlite3", "../database/bddforum.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var email string
		err = db.QueryRow("SELECT email FROM sessions WHERE token = ?", cookie.Value).Scan(&email)
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Database transaction error", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("DELETE FROM users WHERE email = ?", email)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error deleting user", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("DELETE FROM sessions WHERE token = ?", cookie.Value)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error deleting session", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			http.Error(w, "Transaction commit error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Account deleted successfully"))
	}).Methods("POST")

	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err == nil {
			db, err := sql.Open("sqlite3", "../database/bddforum.db")
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

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}).Methods("POST")

}

func ChangePasswordHandlers(router *mux.Router) {
	router.HandleFunc("/front/password/password.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/password/password.html")
	}).Methods("GET")

	router.HandleFunc("/front/password/password.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/password/password.css")
	}).Methods("GET")

	router.HandleFunc("/front/password/password.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/password/password.js")
	}).Methods("GET")

	router.HandleFunc("/api/change-password", ChangePasswordAPIHandler("bddforum.db")).Methods("POST")
}

func ChangePasswordAPIHandler(dbPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangePasswordRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.Password == "" {
			http.Error(w, `{"error": "Mot de passe invalide"}`, http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, `{"error": "Utilisateur non authentifié"}`, http.StatusUnauthorized)
			return
		}

		db, err := sql.Open("sqlite3", "../database/"+dbPath)
		if err != nil {
			http.Error(w, `{"error": "Erreur de connexion à la base de données"}`, http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var email string
		err = db.QueryRow("SELECT email FROM sessions WHERE token = ?", cookie.Value).Scan(&email)
		if err != nil {
			http.Error(w, `{"error": "Session invalide"}`, http.StatusUnauthorized)
			return
		}

		hashedPassword := hashPassword(req.Password)

		_, err = db.Exec("UPDATE users SET password = ? WHERE email = ?", hashedPassword, email)
		if err != nil {
			http.Error(w, `{"error": "Erreur lors de la mise à jour du mot de passe"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Mot de passe mis à jour avec succès"}`))
	}
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
