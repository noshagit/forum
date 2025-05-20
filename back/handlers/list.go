package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func ListPostHandler(router *mux.Router) {
	router.HandleFunc("/post-list", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/post-list/postlist.html")
	}).Methods("GET")

	router.HandleFunc("/post-list", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../script/postlist.js")
	}).Methods("GET")

	router.HandleFunc("/post-list", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/post-list/postlist.css")
	}).Methods("GET")
}
