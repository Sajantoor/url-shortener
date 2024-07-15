package utils

import (
	"context"
	"sync"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthServer struct {
	mu       sync.Mutex
	status   map[string]healthpb.HealthCheckResponse_ServingStatus
	watchers map[string][]chan healthpb.HealthCheckResponse_ServingStatus
}

func NewHealthServer() *HealthServer {
	return &HealthServer{
		status: make(map[string]healthpb.HealthCheckResponse_ServingStatus),
	}
}

func (h *HealthServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	status, ok := h.status[req.Service]
	if !ok {
		status = healthpb.HealthCheckResponse_UNKNOWN
	}

	return &healthpb.HealthCheckResponse{Status: status}, nil
}

func (h *HealthServer) Watch(req *healthpb.HealthCheckRequest, stream healthpb.Health_WatchServer) error {
	statusChan := make(chan healthpb.HealthCheckResponse_ServingStatus, 1)

	h.mu.Lock()
	h.watchers[req.Service] = append(h.watchers[req.Service], statusChan)
	currentStatus := h.status[req.Service]
	h.mu.Unlock()

	if currentStatus == healthpb.HealthCheckResponse_UNKNOWN {
		currentStatus = healthpb.HealthCheckResponse_NOT_SERVING
	}

	if err := stream.Send(&healthpb.HealthCheckResponse{Status: currentStatus}); err != nil {
		return err
	}

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case status := <-statusChan:
			if err := stream.Send(&healthpb.HealthCheckResponse{Status: status}); err != nil {
				return err
			}
		}
	}
}

func (h *HealthServer) SetStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	h.mu.Lock()
	h.status[service] = status
	watchers := h.watchers[service]
	h.mu.Unlock()

	for _, watcher := range watchers {
		watcher <- status
	}
}
