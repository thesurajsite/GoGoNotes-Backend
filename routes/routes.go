package routes

import (
	"github.com/gorilla/mux"
	"github.com/suraj/GoGoNotes/handlers"
)

// setup configures all the routes for the application
func Setup(r *mux.Router, authHandler *handlers.AuthHandler, noteHandler *handlers.NoteHandler) {
	//Auth routes
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	r.HandleFunc("/notes", noteHandler.GetAllNotes).Methods("GET")
	r.HandleFunc("/notes", noteHandler.CreateNote).Methods("POST")
	r.HandleFunc("/notes/{id}", noteHandler.GetNote).Methods("GET")
	r.HandleFunc("/notes/{id}", noteHandler.UpdateNote).Methods("PUT")
	r.HandleFunc("/notes/{id}", noteHandler.DeleteNote).Methods("DELETE")
}
