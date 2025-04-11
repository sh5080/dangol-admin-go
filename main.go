package main

import (
	"context"
	"log"

	"presigned-url-lambda/pkg/config"
	"presigned-url-lambda/pkg/handler"
	"presigned-url-lambda/pkg/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// 어댑터 함수: 의존성 주입을 통해 핸들러 호출
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 설정 로드
	cfg := config.NewConfig()
	
	// S3 클라이언트 초기화
	s3Client, presignClient, err := config.NewS3Client(ctx, cfg)
	if err != nil {
		log.Printf("S3 클라이언트 초기화 실패: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error": "서버 초기화 중 오류가 발생했습니다"}`,
		}, nil
	}
	
	// 서비스 생성
	svc := service.NewPresignedURLService(cfg, s3Client, presignClient)
	
	// 핸들러 생성 및 요청 처리
	h := handler.NewHandler(cfg, svc)
	return h.HandleRequest(ctx, request)
}

func main() {
	lambda.Start(handleRequest)
}
