package taskservice

import (
	"api-service/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	InitialStatus = 0
)

type TaskService struct {
	AgentUrl        string
	DataProviderUrl string
}

func NewTaskService(agent config.AgentConfig, dataProvider config.DataProviderConfig) *TaskService {
	return &TaskService{
		AgentUrl:        agent.Addr,
		DataProviderUrl: dataProvider.Addr,
	}
}

type Status struct {
	Status    int       `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
type Task struct {
	Id            int      `json:"id"`
	StatusHistory []Status `json:"status_history"`
	TaskApi
}

func (svc TaskService) CreateTask(task TaskApi) (*Task, error) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}

	req, err := http.NewRequest("POST", svc.DataProviderUrl, bytes.NewBuffer(jsonTask))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdTask Task
	err = json.NewDecoder(resp.Body).Decode(&createdTask)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	startSolveTask, err := svc.solveTask(createdTask)
	if err != nil {
		return nil, err
	}
	return startSolveTask, nil
}

func (svc TaskService) solveTask(task Task) (*Task, error) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}

	req, err := http.NewRequest("POST", svc.AgentUrl, bytes.NewBuffer(jsonTask))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	task.StatusHistory = append(task.StatusHistory, Status{
		Status:    InitialStatus,
		Timestamp: time.Now(),
	})

	return &task, nil
}
