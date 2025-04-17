package handler

import (
	"context"
	"encoding/json"
	appCtx "lambda-go/pkg/contexts"
	handler "lambda-go/pkg/handlers"
	"lambda-go/pkg/models"
	"lambda-go/pkg/utils"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)
type AdminHandler struct {
	*handler.Handler
}
// GetRestaurantRequests는 매장 생성 요청 목록을 조회합니다.
func (h *AdminHandler) GetRestaurantRequests(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {	
	resp, err := h.AdminService.GetRestaurantRequests(ctx)
	if err != nil {
		return h.HandleAppError(utils.InternalServerError("매장 생성 요청 목록 조회 중 오류가 발생했습니다", err)), nil
	}
	
	return h.SuccessResponse(http.StatusOK, resp), nil
}

// ProcessRestaurantRequest는 매장 생성 요청을 처리합니다.
func (h *AdminHandler) ProcessRestaurantRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// URL 파라미터에서 요청 ID 추출
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
		return h.HandleAppError(utils.InternalServerError("매장 생성 요청 처리 중 오류가 발생했습니다: " + err.Error())), nil
	}
	
	return h.SuccessResponse(http.StatusOK, result), nil
}