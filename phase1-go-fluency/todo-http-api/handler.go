package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidID       = errors.New("invalid task ID")
	ErrNotFound        = errors.New("task not found")
	ErrAlreadyDone     = errors.New("task already completed")
	ErrInvalidTaskName = errors.New("invalid task name")
	ErrInternal        = errors.New("internal error")
)

type CreateTaskRequest struct {
	Description string `json:"description"`
}

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) ListSingleItem(w http.ResponseWriter, ID string) {
	task, err := h.service.ListSingleTask(ID)
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) List(w http.ResponseWriter) {
	tasks, err := h.service.List()
	if err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendError(w, fmt.Errorf("%w:invalid request body", ErrInvalidTaskName))
		return
	}

	task, err := h.service.Add(req.Description)
	if err != nil {
		SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, ID string) {
	if err := h.service.Delete(ID); err != nil {
		SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("task deleted successfully")
}

func (h *TaskHandler) Check(w http.ResponseWriter, ID string) {
	if err := h.service.Check(ID); err != nil {
		SendError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("task completed successfully")
}

func SendError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidID), errors.Is(err, ErrInvalidTaskName):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrAlreadyDone):
		http.Error(w, err.Error(), http.StatusConflict)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
