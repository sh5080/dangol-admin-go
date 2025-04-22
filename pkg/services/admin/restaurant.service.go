package service

import (
	"context"
	"fmt"

	config "lambda-go/pkg/configs"
	"lambda-go/pkg/models"
	repository "lambda-go/pkg/repositories"
	"lambda-go/pkg/utils"
)

// RestaurantService는 매장 관련 서비스를 제공합니다.
type RestaurantService struct {
	config         *config.Config
	restaurantRepo *repository.RestaurantRepository
}

// NewRestaurantService는 새 RestaurantService 인스턴스를 생성합니다.
func NewRestaurantService(cfg *config.Config, restaurantRepo *repository.RestaurantRepository) *RestaurantService {
	return &RestaurantService{
		config:         cfg,
		restaurantRepo: restaurantRepo,
	}
}

// GetRestaurantRequests는 매장 생성 요청 목록을 조회합니다.
func (s *RestaurantService) GetRestaurantRequests(ctx context.Context) (*models.RestaurantRequestsResponse, error) {
	requests, total, err := s.restaurantRepo.GetRestaurantRequests(ctx)
	if err != nil {
		return nil, utils.InternalServerError("매장 생성 요청 목록 조회 실패", err)
	}

	return &models.RestaurantRequestsResponse{
		Requests: requests,
		Total:    total,
	}, nil
}

// ProcessRestaurantRequest는 매장 생성 요청을 승인하거나 거절합니다.
func (s *RestaurantService) ProcessRestaurantRequest(ctx context.Context, requestID string, payload *models.ProcessRestaurantRequest) (*models.RestaurantRequest, error) {
	// 현재 상태 조회
	currentStatus, err := s.restaurantRepo.GetRestaurantRequestByID(ctx, requestID)
	if err != nil {
		return nil, utils.NotFound("요청을 찾을 수 없습니다", err)
	}

	// 이미 처리된 요청인지 확인
	if currentStatus != models.PENDING {
		return nil, utils.BadRequest(fmt.Sprintf("이미 처리된 요청입니다 (현재 상태: %s)", currentStatus))
	}

	// 요청 처리 및 처리된 객체 반환
	result, err := s.restaurantRepo.ProcessRestaurantRequest(ctx, requestID, payload)
	if err != nil {
		return nil, utils.InternalServerError("매장 생성 요청 처리 실패", err)
	}

	return result, nil
}
