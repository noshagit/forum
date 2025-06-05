package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

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
			return
		}

		db, err := sql.Open("sqlite3", "../database/bddforum.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
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
			}
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(credentials.Password)) == nil {
			sessionToken := uuid.New().String()

			_, err := db.Exec("INSERT INTO sessions (token, email) VALUES (?, ?)", sessionToken, credentials.Email)
			if err != nil {
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

	router.PathPrefix("/front/pp/").Handler(http.StripPrefix("/front/pp/", http.FileServer(http.Dir("../../front/pp/"))))

	router.HandleFunc("/front/register/register.html", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "JSON decoding error", http.StatusBadRequest)
			return
		}

		if user.Pseudo == "" || user.Email == "" || user.Password == "" || user.ConfirmPassword == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			log.Printf("Pseudo: %s, Email: %s, Password: %s, ConfirmPassword: %s", user.Pseudo, user.Email, user.Password, user.ConfirmPassword)
			return
		}

		emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		if matched := regexp.MustCompile(emailRegex).MatchString(user.Email); !matched {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		if user.Password != user.ConfirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		pass := user.Password
		if len(pass) < 8 ||
			!regexp.MustCompile(`[A-Z]`).MatchString(pass) ||
			!regexp.MustCompile(`[a-z]`).MatchString(pass) ||
			!regexp.MustCompile(`[0-9]`).MatchString(pass) ||
			!regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(pass) {
			http.Error(w, "Password must be at least 8 characters and include upper, lower, number, and special character", http.StatusBadRequest)
			return
		}

		hashedPassword := hashPassword(user.Password)

		db, err := sql.Open("sqlite3", "../database/bddforum.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var exists int
		err = db.QueryRow("SELECT COUNT(1) FROM users WHERE email = ?", user.Email).Scan(&exists)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if exists > 0 {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}

		err = db.QueryRow("SELECT COUNT(1) FROM users WHERE username = ?", user.Pseudo).Scan(&exists)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		if exists > 0 {
			http.Error(w, "Username already taken", http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("INSERT INTO users (username, email, password, profile_picture) VALUES (?, ?, ?, ?)")
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		defaultProfilePicture := "../pp/default.png"
		_, err = stmt.Exec(user.Pseudo, user.Email, hashedPassword, defaultProfilePicture)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		sessionToken := uuid.New().String()

		_, err = db.Exec("INSERT INTO sessions (token, email) VALUES (?, ?)", sessionToken, user.Email)
		if err != nil {
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
			Pseudo  string `json:"username"`
			Email   string `json:"email"`
			ID      int    `json:"id"`
			Picture string `json:"profile_picture"`
		}

		row = db.QueryRow("SELECT id, username, email, profile_picture FROM users WHERE email = ?", email)
		err = row.Scan(&profile.ID, &profile.Pseudo, &profile.Email, &profile.Picture)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "No profile found for user", http.StatusNotFound)
			} else {
				http.Error(w, "Error retrieving profile information", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"profile": profile,
		})
	}).Methods("GET")

	router.HandleFunc("/profile/{username}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["username"]

		db, err := sql.Open("sqlite3", "../database/bddforum.db")
		if err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		var profile struct {
			Pseudo  string `json:"username"`
			Email   string `json:"email"`
			ID      int    `json:"id"`
			Picture string `json:"profile_picture"`
		}

		row := db.QueryRow("SELECT id, username, email, profile_picture FROM users WHERE username = ?", username)
		err = row.Scan(&profile.ID, &profile.Pseudo, &profile.Email, &profile.Picture)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Profile not found", http.StatusNotFound)
			} else {
				http.Error(w, "Error retrieving profile information", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"profile": profile,
		})
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

		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")

		if username == "" || email == "" {
			http.Error(w, "Username and Email cannot be empty", http.StatusBadRequest)
			return
		}

		emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		if matched := regexp.MustCompile(emailRegex).MatchString(email); !matched {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
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

		var userCount int
		err = db.QueryRow("SELECT COUNT(1) FROM users WHERE username = ? AND email != ?", username, currentEmail).Scan(&userCount)
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		if userCount > 0 {
			http.Error(w, "Username already taken", http.StatusBadRequest)
			return
		}

		var emailCount int
		err = db.QueryRow("SELECT COUNT(1) FROM users WHERE email = ? AND email != ?", email, currentEmail).Scan(&emailCount)
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		if emailCount > 0 {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("profile_picture")
		var filename string
		if err == nil {
			defer file.Close()

			ext := filepath.Ext(handler.Filename)
			filename = fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
			fileSavePath := "../../front/pp/" + filename

			dst, err := os.Create(fileSavePath)
			if err != nil {
				http.Error(w, "Error saving file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, "Error saving file", http.StatusInternalServerError)
				return
			}

			var oldProfilePicture string
			err = db.QueryRow("SELECT profile_picture FROM users WHERE email = ?", currentEmail).Scan(&oldProfilePicture)
			if err != nil {
				http.Error(w, "Error retrieving old profile picture", http.StatusInternalServerError)
				return
			}

			_, err = db.Exec("UPDATE users SET profile_picture = ? WHERE email = ?", "/front/pp/"+filename, currentEmail)
			if err != nil {
				http.Error(w, "Error updating profile picture", http.StatusInternalServerError)
				return
			}

			if oldProfilePicture != "" && oldProfilePicture != "../pp/default.png" && oldProfilePicture != "/front/pp/default.png" {
				oldPath := "../../front/pp/" + filepath.Base(oldProfilePicture)
				if _, err := os.Stat(oldPath); err == nil {
					os.Remove(oldPath)
				}
			}
		}

		_, err = db.Exec("UPDATE users SET username = ?, email = ? WHERE email = ?", username, email, currentEmail)
		if err != nil {
			http.Error(w, "Error updating profile", http.StatusInternalServerError)
			return
		}

		if email != currentEmail {
			_, err = db.Exec("UPDATE sessions SET email = ? WHERE token = ?", email, cookie.Value)
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

		_, err = tx.Exec("UPDATE comments SET owner_id = 0 WHERE owner_id = (SELECT id FROM users WHERE email = ?)", email)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error updating posts", http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("UPDATE posts SET owner_id = 0 WHERE owner_id = (SELECT id FROM users WHERE email = ?)", email)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error updating posts", http.StatusInternalServerError)
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

		pass := req.Password
		if len(pass) < 8 ||
			!regexp.MustCompile(`[A-Z]`).MatchString(pass) ||
			!regexp.MustCompile(`[a-z]`).MatchString(pass) ||
			!regexp.MustCompile(`[0-9]`).MatchString(pass) ||
			!regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(pass) {
			http.Error(w, `{"error": "Le mot de passe doit contenir au moins 8 caractères, une majuscule, une minuscule, un chiffre et un caractère spécial."}`, http.StatusBadRequest)
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
