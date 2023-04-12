package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type ApiServer struct {
	storer ReminderStorer
}

func NewApiServer(storer ReminderStorer) *ApiServer {
	return &ApiServer{
		storer: storer,
	}
}

func (s *ApiServer) Start(listenAddr string) error {
	router := mux.NewRouter()

	router.HandleFunc("/", s.handleGetReminders).Methods("GET")
	router.HandleFunc("/", s.handleCreateReminder).Methods("POST")
	router.HandleFunc("/{id}", s.handleUpdateReminder).Methods("PUT")
	router.HandleFunc("/{id}", s.handleDeleteReminder).Methods("DELETE")

	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	return http.ListenAndServe(fmt.Sprintf(":%s", listenAddr), handlers.CORS(origins, methods)(router))
}

func (s *ApiServer) handleGetReminders(w http.ResponseWriter, r *http.Request) {
	reminders, err := s.storer.GetReminders()
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, reminders)
}

func (s *ApiServer) handleCreateReminder(w http.ResponseWriter, r *http.Request) {
	var reminderBody Reminder
	err := json.NewDecoder(r.Body).Decode(&reminderBody)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	reminder, err := s.storer.CreateReminder(&reminderBody)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, reminder)
}

func (s *ApiServer) handleUpdateReminder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var reminderBody Reminder
	err := json.NewDecoder(r.Body).Decode(&reminderBody)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	reminder, err := s.storer.UpdateReminder(id, &reminderBody)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, reminder)
}

func (s *ApiServer) handleDeleteReminder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	err := s.storer.DeleteReminder(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func writeJSON(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}
