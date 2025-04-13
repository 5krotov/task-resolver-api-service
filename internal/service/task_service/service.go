package taskservice

import (
	"api-service/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/5krotov/task-resolver-pkg/api/v1"
	entity "github.com/5krotov/task-resolver-pkg/entity/v1"
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

func (svc TaskService) CreateTask(task api.CreateTaskRequest) (*entity.Task, error) {
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

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdTask entity.Task
	err = json.NewDecoder(resp.Body).Decode(&createdTask)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &createdTask, nil
}

func (svc TaskService) GetTaskByID(id int64) (*entity.Task, error) {
	url := fmt.Sprintf("%s/api/v1/task/%v", svc.DataProviderUrl, id)
	req, err := http.NewRequest("GET", url, nil)
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

	var task entity.Task
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &task, nil
}
