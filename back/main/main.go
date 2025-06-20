package main

import (
	"Backend/handlers"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func landingPageHandler(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/index.html")
	}).Methods("GET")

	router.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/index.js")
	}).Methods("GET")

	router.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/index.css")
	}).Methods("GET")

	router.HandleFunc("/front/images/bread.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/bread.png")
	}).Methods("GET")

	router.HandleFunc("/front/images/tofu.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/images/tofu.png")
	}).Methods("GET")

}

func main() {

	router := mux.NewRouter()

	landingPageHandler(router)

	handlers.ListPostHandler(router)
	handlers.DetailedPostHandler(router)
	handlers.CommentsHandler(router)
	handlers.LikesHandler(router)

	handlers.RegisterHandler(router)
	handlers.LoginHandler(router)

	handlers.ProfileHandler(router)
	handlers.ChangePasswordHandlers(router)

	fmt.Println("Server is open on port 8080 : http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
