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
	task.HandleFunc("", s.CreateTask).Methods(http.MethodPost)
	task.HandleFunc("", s.GetTasksByFilter).Methods(http.MethodGet)
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

	if task.Difficulty < 0 || task.Difficulty > 3 {
		http.Error(w, "Invalid task difficulty", http.StatusBadRequest)
		return
	}

	createdTask, err := h.svc.CreateTask(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) GetTasksByFilter(w http.ResponseWriter, r *http.Request) {
	perPageStr := r.URL.Query().Get("per_page")
	pageStr := r.URL.Query().Get("page")

	perPage := 20
	page := 1

	if perPageStr != "" {
		if v, err := strconv.Atoi(perPageStr); err == nil && v >= 0 {
			perPage = v
		}
	}
	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v >= 0 {
			page = v
		}
	}

	res, err := h.svc.GetTasksByFilter(perPage, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
