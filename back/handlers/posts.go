package Backend

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func getDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "") // TODO : DB file path
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AddPost(id int, author, title, content string) {
	db, err := getDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO posts (id, author, title, content) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Database preparation error:", err)
		return
	}
	defer stmt.Close()

	if _, err = stmt.Exec(id, author, title, content); err != nil {
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