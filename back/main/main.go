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
		http.ServeFile(w, r, "../script/index.js")
	}).Methods("GET")

	router.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../../front/index.css")
	}).Methods("GET")
}

func main() {

	router := mux.NewRouter()

	landingPageHandler(router)

	handlers.RegisterHandler(router)
	handlers.LoginHandler(router)

	handlers.ProfileHandler(router)
	handlers.LogoutHandler(router)

	fmt.Println("Server is launch on port 8080 : http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
