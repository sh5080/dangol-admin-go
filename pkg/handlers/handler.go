package handler

import (
	"encoding/json"

	config "lambda-go/pkg/configs"
	"lambda-go/pkg/models"
	service "lambda-go/pkg/services"

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
func (h *Handler) corsResponse(statusCode int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
		},
		Body: body,
	}
}

// successResponse는 성공 응답을 생성합니다.
func (h *Handler) successResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	response := models.APIResponse{
		Status: "success",
		Data:   data,
	}
	
	body, err := json.Marshal(response)
	if err != nil {
		return h.errorResponse(500, "응답 생성 중 오류가 발생했습니다")
	}
	
	return h.corsResponse(statusCode, string(body))
}

// 에러 응답 생성 함수
func (h *Handler) errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	response := models.APIErrorResponse{
		Status:  "error",
		Message: message,
	}
	
	body, err := json.Marshal(response)
	if err != nil {
		// JSON 마샬링 실패 시 단순 오류 메시지 반환
		return h.corsResponse(500, `{"status":"error","message":"응답 생성 중 오류가 발생했습니다"}`)
	}
	
	return h.corsResponse(statusCode, string(body))
}