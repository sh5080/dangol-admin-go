package handler

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// HandleRequest는 Lambda 핸들러 함수입니다.
func (h *Handler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("요청 경로: %s, 메서드: %s", request.Path, request.HTTPMethod)
	
	// 요청 경로에 따른 처리
	if strings.HasPrefix(request.Path, "/admin/") {
		// 어드민 API 요청 처리
		return h.AdminHandleRequest(ctx, request)
	} else if request.Path == "/presigned-url" {
		// Presigned URL 요청 처리
		return h.PresignedURLHandleRequest(ctx, request)
	}
	
	// 지원하지 않는 경로
	return h.errorResponse(404, "요청한 API를 찾을 수 없습니다"), nil
} 