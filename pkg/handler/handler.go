package handler

import (
	"context"
	"encoding/json"
	"log"

	"presigned-url-lambda/pkg/config"
	"presigned-url-lambda/pkg/models"
	"presigned-url-lambda/pkg/service"

	"github.com/aws/aws-lambda-go/events"
)

// Handler는 Lambda 핸들러 구조체입니다.
type Handler struct {
	config  *config.Config
	service *service.PresignedURLService
}

// NewHandler는 새 Handler 인스턴스를 생성합니다.
func NewHandler(cfg *config.Config, svc *service.PresignedURLService) *Handler {
	return &Handler{
		config:  cfg,
		service: svc,
	}
}

// HandleRequest는 Lambda 핸들러 함수입니다.
func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("요청 본문: %s", request.Body)

	// OPTIONS 메서드 처리 (CORS 프리플라이트 요청)
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
			},
			Body: "",
		}, nil
	}

	// 요청 파싱
	var req models.PresignedURLRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		log.Printf("요청 파싱 오류: %v", err)
		return errorResponse(400, "잘못된 요청 형식입니다"), nil
	}

	// 요청 검증 및 전처리
	err = h.service.ValidateAndPreprocessRequest(&req)
	if err != nil {
		log.Printf("요청 검증 오류: %v", err)
		return errorResponse(400, err.Error()), nil
	}

	// Presigned URL 생성
	resp, err := h.service.GeneratePresignedURL(ctx, &req)
	if err != nil {
		log.Printf("Presigned URL 생성 오류: %v", err)
		return errorResponse(500, "Presigned URL 생성 중 오류가 발생했습니다"), nil
	}

	// 응답 마샬링
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return errorResponse(500, "응답 생성 중 오류가 발생했습니다"), nil
	}

	// 성공 응답 반환
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
		},
		Body: string(jsonResp),
	}, nil
}

// 에러 응답 생성 함수
func errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(map[string]string{
		"error": message,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
			"Access-Control-Allow-Methods": "GET,POST,OPTIONS",
		},
		Body: string(body),
	}
} 