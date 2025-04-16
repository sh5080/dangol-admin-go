package handler

import (
	config "lambda-go/pkg/configs"
	service "lambda-go/pkg/services"
	"lambda-go/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
)

// Handler는 Lambda 핸들러 구조체입니다.
type Handler struct {
	config       *config.Config
	s3Service    *service.PresignedURL
	adminService *service.Admin
}

// NewHandler는 새 Handler 인스턴스를 생성합니다.
func NewHandler(cfg *config.Config, s3Svc *service.PresignedURL, adminSvc *service.Admin) *Handler {
	return &Handler{
		config:       cfg,
		s3Service:    s3Svc,
		adminService: adminSvc,
	}
}

// corsResponse는 CORS 헤더가 포함된 응답을 생성합니다.
func (h *Handler) corsResponse(response events.APIGatewayProxyResponse) events.APIGatewayProxyResponse {
	// 기존 헤더 유지하면서 CORS 헤더 추가
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}
	
	response.Headers["Access-Control-Allow-Origin"] = "*"
	response.Headers["Access-Control-Allow-Headers"] = "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"
	response.Headers["Access-Control-Allow-Methods"] = "GET,POST,OPTIONS"
	
	return response
}

// successResponse는 성공 응답을 생성합니다.
func (h *Handler) successResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	response, err := utils.Success(statusCode, data)
	if err != nil {
		return h.errorResponse(500, "응답 생성 중 오류가 발생했습니다")
	}
	
	return h.corsResponse(response)
}

// errorResponse는 에러 응답을 생성합니다.
func (h *Handler) errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	response, err := utils.Error(statusCode, message)
	if err != nil {
		// JSON 마샬링에 실패한 경우에 대한 기본 응답
		defaultResponse := events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"status":"error","message":"응답 생성 중 오류가 발생했습니다"}`,
		}
		return h.corsResponse(defaultResponse)
	}
	
	return h.corsResponse(response)
}

// handleAppError는 AppError를 적절한 API 응답으로 변환합니다.
func (h *Handler) handleAppError(appErr *utils.AppError) events.APIGatewayProxyResponse {
	return h.errorResponse(appErr.StatusCode, appErr.Message)
}

