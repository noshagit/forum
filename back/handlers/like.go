package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func getDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "../database/bddforum.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func LikesHandler(router *mux.Router) {
	router.HandleFunc("/api/like", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../api/like.js")
	}).Methods("GET")

	router.HandleFunc("/posts/like", LikePostHandler).Methods("POST")
	router.HandleFunc("/posts/unlike", UnlikePostHandler).Methods("POST")
	router.HandleFunc("/posts/{post_id}/like_count", GetPostLikes).Methods("GET")

	router.HandleFunc("/comments/like", LikeCommentHandler).Methods("POST")
	router.HandleFunc("/comments/unlike", UnlikeCommentHandler).Methods("POST")
	router.HandleFunc("/comments/{comment_id}/like_count", GetCommentLikes).Methods("GET")
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID int `json:"userId"`
		PostID int `json:"postId"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, postID := requestBody.UserID, requestBody.PostID

	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT OR IGNORE INTO post_likes (post_id, owner_id) VALUES (?, ?)`)
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(postID, userID); err != nil {
		log.Println("Database insertion error:", err)
		return
	}

	stmt, err = db.Prepare(`UPDATE posts SET likes = (SELECT COUNT(*) FROM post_likes WHERE post_id = ?) WHERE id = ?`)
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	if _, err = stmt.Exec(postID, postID); err != nil {
		log.Println("Database update error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UnlikePostHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID int `json:"userId"`
		PostID int `json:"postId"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, postID := requestBody.UserID, requestBody.PostID

	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM post_likes WHERE post_id = ? AND owner_id = ?`)
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(postID, userID); err != nil {
		log.Println("Database deletion error:", err)
		return
	}

	stmt, err = db.Prepare(`UPDATE posts SET likes = (SELECT COUNT(*) FROM post_likes WHERE post_id = ?) WHERE id = ?`)
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	if _, err = stmt.Exec(postID, postID); err != nil {
		log.Println("Database update error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetPostLikes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["post_id"]

	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	var count int
	err = db.QueryRow(`SELECT likes FROM posts WHERE id = ?`, postID).Scan(&count)
	if err != nil {
		count = 0
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"likes": count})
}

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID    int `json:"userId"`
		CommentID int `json:"commentId"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, commentID := requestBody.UserID, requestBody.CommentID

	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(`INSERT OR IGNORE INTO comment_likes (comment_id, owner_id) VALUES (?, ?)`)
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(commentID, userID); err != nil {
		log.Println("Database insertion error:", err)
		return
	}

	stmt, err = db.Prepare(`UPDATE comments SET likes = (SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?) WHERE id = ?`)
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	if _, err = stmt.Exec(commentID, commentID); err != nil {
		log.Println("Database update error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UnlikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		UserID    int `json:"userId"`
		CommentID int `json:"commentId"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, commentID := requestBody.UserID, requestBody.CommentID

	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM comment_likes WHERE comment_id = ? AND owner_id = ?`)
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err := stmt.Exec(commentID, userID); err != nil {
		log.Println("Database deletion error:", err)
		return
	}

	stmt, err = db.Prepare(`UPDATE comments SET likes = (SELECT COUNT(*) FROM comment_likes WHERE comment_id = ?) WHERE id = ?`)
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	if _, err = stmt.Exec(commentID, commentID); err != nil {
		log.Println("Database update error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetCommentLikes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["comment_id"]

	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	var count int
	err = db.QueryRow(`SELECT likes FROM comments WHERE id = ?`, commentID).Scan(&count)
	if err != nil {
		count = 0
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"likes": count})
}
