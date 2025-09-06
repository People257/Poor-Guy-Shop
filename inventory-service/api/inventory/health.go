package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HealthCheckResponse 健康检查响应
type HealthCheckResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
	Errors   map[string]string `json:"errors,omitempty"`
}

// HealthCheckHandler HTTP健康检查处理器
func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	response := HealthCheckResponse{
		Status:   "healthy",
		Services: make(map[string]string),
		Errors:   make(map[string]string),
	}

	// 检查业务服务
	if s.businessService != nil {
		healthResults := s.businessService.HealthCheck(ctx)
		for service, err := range healthResults {
			if err != nil {
				response.Services[service] = "unhealthy"
				response.Errors[service] = err.Error()
				response.Status = "unhealthy"
			} else {
				response.Services[service] = "healthy"
			}
		}
	}

	// 设置HTTP状态码
	statusCode := http.StatusOK
	if response.Status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

// ReadinessCheck gRPC就绪检查
func (s *Server) ReadinessCheck(ctx context.Context) error {
	if s.businessService == nil {
		return fmt.Errorf("business service is not initialized")
	}

	healthResults := s.businessService.HealthCheck(ctx)
	for service, err := range healthResults {
		if err != nil {
			return fmt.Errorf("service %s is unhealthy: %w", service, err)
		}
	}

	return nil
}

// LivenessCheck gRPC存活检查
func (s *Server) LivenessCheck(ctx context.Context) error {
	// 简单的存活检查，确保服务器可以响应
	if s.inventoryApp == nil {
		return fmt.Errorf("inventory application is not initialized")
	}

	return nil
}

