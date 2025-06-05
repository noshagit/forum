package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func CommentsHandler(router *mux.Router) {
	router.HandleFunc("/api/get_avatar/{Author}", GetProfilePicture).Methods("GET")

	router.HandleFunc("/api/add-comment", AddCommentHandler).Methods("POST")

	router.HandleFunc("/api/delete-comment", deleteCommentHandler).Methods("DELETE")

	router.HandleFunc("/api/edit-comment", ModifyCommentHandler).Methods("PUT")
}

func AddComment(postID, ownerID int, content string) {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO comments (post_id, owner_id, content) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(postID, ownerID, content); err != nil {
		log.Println("Database insertion error:", err)
		return
	}
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value

	db, err := getDB()
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var req struct {
		PostID  string `json:"postId"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	postID, err := strconv.Atoi(req.PostID)
	if err != nil {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		log.Println("Erreur de conversion d'ID de post:", err)
		return
	}

	var email string
	err = db.QueryRow("SELECT email FROM sessions WHERE token = ?", sessionToken).Scan(&email)
	if err != nil {
		http.Error(w, "Session invalide", http.StatusUnauthorized)
		return
	}

	var authorID int
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&authorID)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusInternalServerError)
		return
	}

	AddComment(postID, authorID, req.Content)
	w.WriteHeader(http.StatusOK)
}

func DeleteComment(id int) {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM comments WHERE id = ?")
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id); err != nil {
		log.Println("Database deletion error:", err)
		return
	}
}

func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}
	sessionToken := cookie.Value

	db, err := getDB()
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var email string
	err = db.QueryRow("SELECT email FROM sessions WHERE token = ?", sessionToken).Scan(&email)
	if err != nil {
		http.Error(w, "Session invalide", http.StatusUnauthorized)
		return
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusInternalServerError)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	var ownerID int
	err = db.QueryRow("SELECT owner_id FROM comments WHERE id = ?", id).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Commentaire introuvable", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Vous n'êtes pas propriétaire de ce commentaire", http.StatusForbidden)
		return
	}

	DeleteComment(id)
	w.WriteHeader(http.StatusOK)
}

func ModifyComment(id int, newContent string) error {
	db, err := getDB()
	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE comments SET content = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("preparation error: %v", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(newContent, id)
	if err != nil {
		return fmt.Errorf("execution error: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération des lignes modifiées: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("aucune ligne modifiée. L’ID %d existe-t-il ?", id)
	}

	return nil
}

func ModifyCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ID      int    `json:"id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Non connecté", http.StatusUnauthorized)
		return
	}
	sessionToken := cookie.Value

	db, err := getDB()
	if err != nil {
		http.Error(w, "Erreur DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var email string
	err = db.QueryRow("SELECT email FROM sessions WHERE token = ?", sessionToken).Scan(&email)
	if err != nil {
		http.Error(w, "Session invalide", http.StatusUnauthorized)
		return
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&userID)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusInternalServerError)
		return
	}

	var ownerID int
	err = db.QueryRow("SELECT owner_id FROM comments WHERE id = ?", data.ID).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Vous n'êtes pas propriétaire de ce post", http.StatusForbidden)
		return
	}

	if err := ModifyComment(data.ID, data.Content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetProfilePicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["Author"]

	db, err := getDB()
	if err != nil {
		fmt.Println("Database connection error:", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var avatarPath string
	err = db.QueryRow("SELECT profile_picture FROM users WHERE username = ?", username).Scan(&avatarPath)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	var path string
	if strings.Contains(avatarPath, "default.png") {
		path = "../../front/pp/default.png"
	} else {
		path = "../.." + avatarPath
	}
	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, path)
}
