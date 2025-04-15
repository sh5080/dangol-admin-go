package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"lambda-go/pkg/models"
	"lambda-go/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
)

// AdminHandleRequest는 어드민 API 요청을 처리합니다.
func (h *Handler) AdminHandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := request.Path
	method := request.HTTPMethod

	// OPTIONS 메서드 처리 (CORS 프리플라이트 요청)
	if method == "OPTIONS" {
		return h.corsResponse(http.StatusOK, ""), nil
	}

	// 요청 경로에 따라 처리
	switch {
	case path == "/admin/restaurant/request" && method == "GET":
		return h.getRestaurantRequests(ctx)
	case strings.HasPrefix(path, "/admin/restaurant/request/") && method == "POST":
		parts := strings.Split(path, "/")
		if len(parts) < 5 || parts[4] != "process" {
			return h.errorResponse(http.StatusBadRequest, "잘못된 API 경로입니다"), nil
		}

		requestID, err := strconv.Atoi(parts[3])
		if err != nil {
			return h.errorResponse(http.StatusBadRequest, "유효하지 않은 요청 ID입니다"), nil
		}

		return h.processRestaurantRequest(ctx, requestID, request.Body)
	default:
		return h.errorResponse(http.StatusNotFound, "요청한 API를 찾을 수 없습니다"), nil
	}
}

// getRestaurantRequests는 매장 생성 요청 목록을 조회합니다.
func (h *Handler) getRestaurantRequests(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	resp, err := h.adminService.GetRestaurantRequests(ctx)
	if err != nil {
		log.Printf("매장 생성 요청 목록 조회 오류: %v", err)
		return h.errorResponse(http.StatusInternalServerError, "매장 생성 요청 목록 조회 중 오류가 발생했습니다"), nil
	}

	return h.successResponse(http.StatusOK, resp), nil
}

// processRestaurantRequest는 매장 생성 요청을 처리합니다.
func (h *Handler) processRestaurantRequest(ctx context.Context, requestID int, body string) (events.APIGatewayProxyResponse, error) {
	var payload models.ProcessRestaurantRequest

	err := json.Unmarshal([]byte(body), &payload)
	if err != nil {
		log.Printf("요청 파싱 오류: %v", err)
		return h.errorResponse(http.StatusBadRequest, "잘못된 요청 형식입니다"), nil
	}
	if err := utils.Validate(&payload); err != nil {
		log.Printf("유효성 검증 오류: %v", err)
		return h.errorResponse(http.StatusBadRequest, err.Error()), nil
	}

	result, err := h.adminService.ProcessRestaurantRequest(ctx, requestID, &payload)
	if err != nil {
		log.Printf("매장 생성 요청 처리 오류: %v", err)
		return h.errorResponse(http.StatusInternalServerError, fmt.Sprintf("매장 생성 요청 처리 중 오류가 발생했습니다: %s", err.Error())), nil
	}


	return h.successResponse(http.StatusOK, result), nil
}