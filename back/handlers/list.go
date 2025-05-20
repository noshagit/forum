package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ListPostHandler(router *mux.Router) {
	router.HandleFunc("/post-list/postlist.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/post-list/postlist.html")
	}).Methods("GET")

	router.HandleFunc("/post-list/post-list.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/post-list/postlist.js")
	}).Methods("GET")

	router.HandleFunc("/post-list/post-list.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/post-list/postlist.css")
	}).Methods("GET")
}
