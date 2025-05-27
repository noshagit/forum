package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

type Comment struct {
	ID        int
	PostID    int
	OwnerID   int
	Author    string
	Content   string
	CreatedAt string
}

type Post struct {
	ID        int
	OwnerID   int
	Title     string
	Content   string
	Likes     int
	Themes    string
	CreatedAt string
	Author    string
}

func ListPostHandler(router *mux.Router) {
	router.HandleFunc("/front/post-list/postlist.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/post-list/postlist.html")
	}).Methods("GET")

	router.HandleFunc("/front/post-list/postlist.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/post-list/postlist.css")
	}).Methods("GET")

	router.HandleFunc("/front/post-list/postlist.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/post-list/postlist.js")
	}).Methods("GET")

	router.HandleFunc("/front/images/like.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/like.png")
	}).Methods("GET")

	router.HandleFunc("/front/images/share.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/share.png")
	}).Methods("GET")

	router.HandleFunc("/front/images/comment.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/comment.png")
	}).Methods("GET")

	router.HandleFunc("/front/images/logo-minecraft.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/logo-minecraft.png")
	}).Methods("GET")

	router.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		posts := GetPosts()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
		}
	}).Methods("GET")

	router.HandleFunc("/api/add-post", AddPostHandler).Methods("POST")

	router.HandleFunc("/api/modify-post", ModifyPostHandler).Methods("PUT")

	router.HandleFunc("/api/delete-post", deletePostHandler).Methods("DELETE")
}

func GetUserIDFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}

	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return 0, err
	}
	defer db.Close()

	sessionToken := cookie.Value

	var userID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", sessionToken).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func GetPosts() []Post {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return nil
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, owner_id, title, content, likes, themes, created_at FROM posts")
	if err != nil {
		log.Println("Database query error:", err)
		return nil
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.OwnerID, &post.Title, &post.Content, &post.Likes, &post.Themes, &post.CreatedAt); err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		posts = append(posts, post)
	}

	return posts
}

func AddPost(author int, title, content, themes string) {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO posts (owner_id, title, content, themes) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(author, title, content, themes); err != nil {
		log.Println("Database insertion error:", err)
		return
	}
}

func AddPostHandler(w http.ResponseWriter, r *http.Request) {
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
	log.Println("Session token:", sessionToken)

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

	var authorID int
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&authorID)
	if err != nil {
		http.Error(w, "Utilisateur introuvable", http.StatusInternalServerError)
		return
	}

	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Themes  string `json:"themes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON invalide", http.StatusBadRequest)
		return
	}

	AddPost(authorID, req.Title, req.Content, req.Themes)
	w.WriteHeader(http.StatusOK)
}

func DeletePost(id int) {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM posts WHERE id = ?")
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

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
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

	// Ouvrir la base de données
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
	err = db.QueryRow("SELECT owner_id FROM posts WHERE id = ?", id).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Vous n'êtes pas propriétaire de ce post", http.StatusForbidden)
		return
	}

	DeletePost(id)
	w.WriteHeader(http.StatusOK)
}

func ModifyPost(id int, title string, newContent string, themes string) error {
	log.Printf("Modification en cours pour ID=%d, titre='%s', contenu='%s', thème='%s'", id, title, newContent, themes)
	db, err := getDB()
	if err != nil {
		return fmt.Errorf("database connection error: %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE posts SET title = ?, content = ?, themes = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("preparation error: %v", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(title, newContent, themes, id)
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

func ModifyPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Theme   string `json:"theme"`
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
	err = db.QueryRow("SELECT owner_id FROM posts WHERE id = ?", data.ID).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Post introuvable", http.StatusNotFound)
		return
	}

	if ownerID != userID {
		http.Error(w, "Vous n'êtes pas propriétaire de ce post", http.StatusForbidden)
		return
	}

	if err := ModifyPost(data.ID, data.Title, data.Content, data.Theme); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetComments(postID int) []Comment {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return nil
	}
	defer db.Close()

	query := ` SELECT c.id, c.post_id, c.owner_id, u.username, c.content, c.created_at FROM comments c LEFT JOIN users u ON c.owner_id = u.id WHERE c.post_id = ? `
	rows, err := db.Query(query, postID)
	if err != nil {
		log.Println("Database query error:", err)
		return nil
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.OwnerID, &comment.Author, &comment.Content, &comment.CreatedAt)
		if err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		comments = append(comments, comment)
	}

	return comments
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postIDStr := vars["postID"]

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	comments := GetComments(postID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
	}
}

func DetailedPostHandler(router *mux.Router) {
	router.HandleFunc("/front/comments/comments.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/comments/comments.html")
	}).Methods("GET")

	router.HandleFunc("/front/comments/comments.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/comments/comments.css")
	}).Methods("GET")

	router.HandleFunc("/front/comments/comments.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/comments/comments.js")
	}).Methods("GET")

	router.HandleFunc("/api/comments/{postID}", GetCommentsHandler).Methods("GET")

	router.HandleFunc("/api/post/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID invalide", http.StatusBadRequest)
			return
		}

		db, err := sql.Open("sqlite3", "../database/bddforum.db")
		if err != nil {
			http.Error(w, "Erreur base de données", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		query := ` SELECT p.id, p.title, p.content, p.themes, p.likes, p.created_at, u.username FROM posts p LEFT JOIN users u ON p.owner_id = u.id WHERE p.id = ?`

		row := db.QueryRow(query, id)

		var post Post
		err = row.Scan(&post.ID, &post.Title, &post.Content, &post.Themes, &post.Likes, &post.CreatedAt, &post.Author)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Post non trouvé", http.StatusNotFound)
				return
			}
			http.Error(w, "Erreur lors de la lecture du post", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}).Methods("GET")
}
