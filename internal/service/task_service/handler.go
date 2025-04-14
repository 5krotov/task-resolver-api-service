package taskservice

import (
	"encoding/json"
	"net/http"
	"strconv"

	api "github.com/5krotov/task-resolver-pkg/api/v1"
	mux "github.com/gorilla/mux"
)

type TaskHandler struct {
	svc *TaskService
}

func NewTaskHandler(svc *TaskService) *TaskHandler {
	return &TaskHandler{
		svc: svc,
	}
}

func (s *TaskHandler) Register(mux *mux.Router) {
	task := mux.PathPrefix("/api/v1/task").Subrouter()
	task.HandleFunc("", s.CreateTask)
	task.HandleFunc("/{id}", s.GetTaskByID)
}

func (h TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task api.CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Name == "" {
		http.Error(w, "Task name is required", http.StatusBadRequest)
		return
	}

	if task.Difficulty < 0 || task.Difficulty > 10 {
		http.Error(w, "Invalid task difficulty", http.StatusBadRequest)
		return
	}

	createdTask, err := h.svc.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func (h TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.svc.GetTaskByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
