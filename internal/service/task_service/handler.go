package taskservice

import (
	"encoding/json"
	"net/http"
)

type TaskHandler struct {
	svc *TaskService
}

func NewTaskHandler(svc *TaskService) *TaskHandler {
	return &TaskHandler{
		svc: svc,
	}
}

func (s *TaskHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/task", s.CreateTask)
}

type TaskApi struct {
	Name       string `json:"name"`
	Difficulty int    `json:"difficulty"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task TaskApi
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
