package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gorilla/mux"
)

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
}

type Post struct {
	ID        int
	OwnerID   int
	Title     string
	Content   string
	Likes     int
	Themes    string
	CreatedAt string
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

func AddPost(id int, author, title, content string) {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO posts (id, owner_id, title, content, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id, author, title, content, time.Now().Format("2006-01-02 15:04")); err != nil {
		log.Println("Database insertion error:", err)
		return
	}
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

func ModifyPost(id int, newContent string) {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE posts SET content = ? WHERE id = ?")
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

type Comment struct {
	ID        int
	PostID    int
	OwnerID   int
	Content   string
	CreatedAt string
}

func GetComments(postID int) []Comment {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return nil
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, post_id, owner_id, content, created_at FROM comments WHERE post_id = ?", postID)
	if err != nil {
		log.Println("Database query error:", err)
		return nil
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.OwnerID, &comment.Content, &comment.CreatedAt); err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		comments = append(comments, comment)
	}

	return comments
}
