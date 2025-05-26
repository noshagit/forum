package handlers

import (
	"net/http"

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
}
