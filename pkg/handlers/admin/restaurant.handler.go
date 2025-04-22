package handler

import (
	"context"
	"encoding/json"
	appCtx "lambda-go/pkg/contexts"
	handler "lambda-go/pkg/handlers"
	"lambda-go/pkg/models"
	"lambda-go/pkg/utils"
	"net/http"

	dto "lambda-go/pkg/models/dtos"

	"github.com/aws/aws-lambda-go/events"
)

type AdminHandler struct {
	*handler.Handler
}

// GetRestaurantRequests는 매장 생성 요청 목록을 조회합니다.
func (h *AdminHandler) GetRestaurantRequests(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	query := dto.RestaurantRequestQuery{}

	query.Page, query.PageSize = appCtx.ParsePaginationParams(request)

	statusStr := appCtx.GetStringParam(request, "status", "")
	if statusStr != "" {
		status := models.RestaurantRequestStatus(statusStr)
		if status == models.PENDING || status == models.APPROVED || status == models.REJECTED {
			query.Status = &status
		}
	}

	resp, err := h.AdminService.GetRestaurantRequests(ctx, query)
	if err != nil {
		return h.HandleAppError(err), nil
	}

	return h.SuccessResponse(http.StatusOK, resp), nil
}

// ProcessRestaurantRequest는 매장 생성 요청을 처리합니다.
func (h *AdminHandler) ProcessRestaurantRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestID := appCtx.GetParam(ctx, "id")
	if requestID == "" {
		return h.HandleAppError(utils.BadRequest("유효하지 않은 요청 ID입니다")), nil
	}

	var payload models.ProcessRestaurantRequest

	err := json.Unmarshal([]byte(request.Body), &payload)
	if err != nil {
		return h.HandleAppError(utils.BadRequest("잘못된 요청 형식입니다: " + err.Error())), nil
	}

	if err := utils.Validate(&payload); err != nil {
		return h.HandleAppError(utils.BadRequest(err.Error())), nil
	}

	result, err := h.AdminService.ProcessRestaurantRequest(ctx, requestID, &payload)
	if err != nil {
		return h.HandleAppError(err), nil
	}

	return h.SuccessResponse(http.StatusOK, result), nil
}
