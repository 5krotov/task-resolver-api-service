package taskservice

import (
	"api-service/internal/config"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"

	pb "github.com/5krotov/task-resolver-pkg/grpc-api/v1"
	"github.com/5krotov/task-resolver-pkg/rest-api/v1/api"
	"github.com/5krotov/task-resolver-pkg/rest-api/v1/entity"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskService struct {
	AgentConn          grpc.ClientConnInterface
	DataProviderConn   grpc.ClientConnInterface
	AgentClient        pb.AgentServiceClient
	DataProviderClient pb.DataProviderServiceClient
	logger             *zap.Logger
}

func NewTaskService(agent config.AgentConfig, dataProvider config.DataProviderConfig, grpcClient config.GRPCClientConfig, logger *zap.Logger) *TaskService {
	var err error

	var clientCert tls.Certificate
	if len(grpcClient.ClientCert) != 0 && len(grpcClient.ClientKey) != 0 {
		clientCert, err = tls.LoadX509KeyPair(grpcClient.ClientCert, grpcClient.ClientKey)
		if err != nil {
			log.Fatalf("Failed to load client cert: %v", err)
		}
	}

	var caCert []byte
	if len(grpcClient.CaCert) != 0 {
		caCert, err = os.ReadFile(grpcClient.CaCert)
		if err != nil {
			log.Fatalf("Failed to read CA cert: %v", err)
		}
	}

	var agentConn grpc.ClientConnInterface
	if agent.UseTLS {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			logger.Fatal("failed to add CA certificate for agent service")
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs:      caCertPool,
			ServerName:   agent.GrpcServerName,
		}

		agentConn, err = grpc.NewClient(
			agent.Addr,
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		)
		if err != nil {
			logger.Fatal(fmt.Sprintf("failed to connect to %v, error: %v", agent.Addr, err))
		}
	} else {
		agentConn, err = grpc.NewClient(agent.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal(fmt.Sprintf("failed to connect to %v, error: %v", agent.Addr, err))
		}
	}

	agentServiceClient := pb.NewAgentServiceClient(agentConn)

	var dataProviderConn grpc.ClientConnInterface
	if dataProvider.UseTLS {
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			logger.Fatal("failed to add CA certificate for dataProvider service")
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs:      caCertPool,
			ServerName:   dataProvider.GrpcServerName,
		}

		dataProviderConn, err = grpc.NewClient(
			dataProvider.Addr,
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		)
		if err != nil {
			logger.Fatal(fmt.Sprintf("failed to connect to %v, error: %v", dataProvider.Addr, err))
		}
	} else {
		dataProviderConn, err = grpc.NewClient(dataProvider.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Fatal(fmt.Sprintf("failed to connect to %v, error: %v", dataProvider.Addr, err))
		}
	}

	dataProviderServiceClient := pb.NewDataProviderServiceClient(dataProviderConn)

	return &TaskService{
		AgentConn:          agentConn,
		DataProviderConn:   dataProviderConn,
		AgentClient:        agentServiceClient,
		DataProviderClient: dataProviderServiceClient,
		logger:             logger,
	}
}

func (svc TaskService) CreateTask(ctx context.Context, task api.CreateTaskRequest) (*entity.Task, error) {
	svc.logger.Info("Creating task", zap.Any("task", task))

	req := &pb.CreateTaskRequest{
		Name:       task.Name,
		Difficulty: pb.Difficulty(task.Difficulty),
	}
	resp, err := svc.AgentClient.CreateTask(ctx, req)
	if err != nil {
		svc.logger.Error("Error from agent", zap.String("error", err.Error()))
		return nil, fmt.Errorf("Error from agent: %s", err.Error())
	}

	svc.logger.Info("Task created successfully", zap.Any("task", resp.Task))

	return mapPbTaskToEntityTask(resp.Task), nil
}

func mapEntityDoTaskRequestToPbDoTaskRequest(task *pb.Task) *entity.Task {
	statusHistory := []entity.Status{}
	for _, s := range task.GetStatusHistory() {
		statusHistory = append(statusHistory, entity.Status{
			Status:    int(s.GetStatus()),
			Timestamp: s.GetTimestamp().AsTime(),
		})
	}
	return &entity.Task{
		Id:            task.GetId(),
		Name:          task.GetName(),
		Difficulty:    int(task.GetDifficulty()),
		StatusHistory: statusHistory,
	}
}

func mapPbTaskToEntityTask(task *pb.Task) *entity.Task {
	statusHistory := []entity.Status{}
	for _, s := range task.GetStatusHistory() {
		statusHistory = append(statusHistory, entity.Status{
			Status:    int(s.GetStatus()),
			Timestamp: s.GetTimestamp().AsTime(),
		})
	}
	return &entity.Task{
		Id:            task.GetId(),
		Name:          task.GetName(),
		Difficulty:    int(task.GetDifficulty()),
		StatusHistory: statusHistory,
	}
}

func (svc TaskService) GetTaskByID(ctx context.Context, id int64) (*entity.Task, error) {
	svc.logger.Info("Getting task by ID", zap.Int64("id", id))

	req := &pb.GetTaskRequest{
		Id: id,
	}
	resp, err := svc.DataProviderClient.GetTask(ctx, req)
	if err != nil {
		svc.logger.Error("Error from data provider", zap.String("error", err.Error()))
		return nil, fmt.Errorf("Error from data provider: %s", err.Error())
	}

	return mapPbTaskToEntityTask(resp.Task), nil
}

func (svc TaskService) GetTasksByFilter(ctx context.Context, per_page, page int) (*api.SearchTaskResponse, error) {
	svc.logger.Info("Getting tasks by filter", zap.Int("per_page", per_page), zap.Int("page", page))

	req := &pb.SearchTaskRequest{
		Page:    int64(page),
		PerPage: int64(per_page),
	}
	resp, err := svc.DataProviderClient.SearchTask(ctx, req)
	if err != nil {
		svc.logger.Error("Error from data provider", zap.String("error", err.Error()))
		return nil, fmt.Errorf("Error from data provider: %s", err.Error())
	}

	resultResp := api.SearchTaskResponse{
		Pages: int(resp.Pages),
		Tasks: make([]entity.Task, 0),
	}
	for _, st := range resp.Tasks {
		resultResp.Tasks = append(resultResp.Tasks, *mapPbTaskToEntityTask(st))
	}
	return &resultResp, nil
}
