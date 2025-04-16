package taskservice

import (
	"api-service/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	api "github.com/5krotov/task-resolver-pkg/api/v1"
	entity "github.com/5krotov/task-resolver-pkg/entity/v1"
	"go.uber.org/zap"
)

type TaskService struct {
	AgentUrl        string
	DataProviderUrl string
	logger          *zap.Logger
}

func NewTaskService(agent config.AgentConfig, dataProvider config.DataProviderConfig, logger *zap.Logger) *TaskService {
	return &TaskService{
		AgentUrl:        agent.Addr,
		DataProviderUrl: dataProvider.Addr,
		logger:          logger,
	}
}

func (svc TaskService) CreateTask(task api.CreateTaskRequest) (*entity.Task, error) {
	svc.logger.Info("Creating task", zap.Any("task", task))

	jsonTask, err := json.Marshal(task)
	if err != nil {
		svc.logger.Error("Failed to marshal task", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/task", svc.AgentUrl)
	svc.logger.Info("Sending request to agent", zap.String("url", url))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonTask))
	if err != nil {
		svc.logger.Error("Failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		svc.logger.Error("Failed to send request", zap.Error(err))
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		svc.logger.Error("Unexpected status code from agent", zap.Int("statusCode", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code of agent: %d", resp.StatusCode)
	}

	var createdTask entity.Task
	err = json.NewDecoder(resp.Body).Decode(&createdTask)
	if err != nil {
		svc.logger.Error("Failed to unmarshal response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	svc.logger.Info("Task created successfully", zap.Any("task", createdTask))
	return &createdTask, nil
}

func (svc TaskService) GetTaskByID(id int64) (*entity.Task, error) {
	svc.logger.Info("Getting task by ID", zap.Int64("id", id))

	url := fmt.Sprintf("%s/api/v1/task/%v", svc.DataProviderUrl, id)
	svc.logger.Info("Sending request to data provider", zap.String("url", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		svc.logger.Error("Failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		svc.logger.Error("Failed to send request", zap.Error(err))
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		svc.logger.Error("Unexpected status code from data provider", zap.Int("statusCode", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code of dataprovider: %d", resp.StatusCode)
	}

	var task entity.Task
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		svc.logger.Error("Failed to unmarshal response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	svc.logger.Info("Task retrieved successfully", zap.Any("task", task))
	return &task, nil
}

func (svc TaskService) GetTasksByFilter(per_page, page int) (*api.SearchTaskResponse, error) {
	svc.logger.Info("Getting tasks by filter", zap.Int("per_page", per_page), zap.Int("page", page))

	url := fmt.Sprintf("%s/api/v1/task?per_page=%d&page=%d", svc.DataProviderUrl, per_page, page)
	svc.logger.Info("Sending request to data provider", zap.String("url", url))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		svc.logger.Error("Failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		svc.logger.Error("Failed to send request", zap.Error(err))
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		svc.logger.Error("Unexpected status code from data provider", zap.Int("statusCode", resp.StatusCode))
		return nil, fmt.Errorf("unexpected status code of dataprovider: %d", resp.StatusCode)
	}

	var result api.SearchTaskResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		svc.logger.Error("Failed to unmarshal response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	svc.logger.Info("Tasks retrieved successfully", zap.Any("result", result))
	return &result, nil
}
