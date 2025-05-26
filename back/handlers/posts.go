package handlers

import (
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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
