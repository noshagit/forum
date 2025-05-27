package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CommentsHandler(router *mux.Router) {
	router.HandleFunc("/api/get_avatar/{userID}", GetProfilePicture).Methods("GET")

}

func AddComment(id, postID int, author, content string) {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO comments (id, post_id, owner_id, content, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id, postID, author, content, time.Now().Format("2006-01-02 15:04")); err != nil {
		log.Println("Database insertion error:", err)
		return
	}

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

func ModifyComment(id int, newContent string) {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE comments SET content = ? WHERE id = ?")
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(newContent, id); err != nil {
		log.Println("Database insertion error:", err)
		return
	}
}

func GetProfilePicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["Author"]

	db, err := getDB()
	if err != nil {
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

	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeFile(w, r, avatarPath)
	if err != nil {
		log.Println("Error serving profile picture:", err)
		http.Error(w, "Error serving profile picture", http.StatusInternalServerError)
		return
	}
}
