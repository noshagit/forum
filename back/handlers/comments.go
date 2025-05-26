package handlers

import (
	"log"
	"time"
)

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
