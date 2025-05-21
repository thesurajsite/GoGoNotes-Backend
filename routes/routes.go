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

	// Note routes (unprotected)
	noteRouter := r.PathPrefix("/notes").Subrouter()

	noteRouter.HandleFunc("", noteHandler.GetAllNotes).Methods("GET")
	noteRouter.HandleFunc("", noteHandler.CreateNote).Methods("POST")
	noteRouter.HandleFunc("/{id}", noteHandler.GetNote).Methods("GET")
	noteRouter.HandleFunc("/{id}", noteHandler.UpdateNote).Methods("PUT")
	noteRouter.HandleFunc("/{id}", noteHandler.DeleteNote).Methods("DELETE")
}
