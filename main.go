package main

import (
	"context"
	"fmt"
	"log"

	config "lambda-go/pkg/configs"
	handler "lambda-go/pkg/handlers"
	repository "lambda-go/pkg/repositories"
	service "lambda-go/pkg/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v4/pgxpool"
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
	
	// 데이터베이스 초기화 (pgxpool 사용)
	dbCfg := cfg.NewDBConfig()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.Database, dbCfg.SSLMode)
	
	dbPool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		log.Printf("데이터베이스 연결 실패: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error": "데이터베이스 연결 중 오류가 발생했습니다"}`,
		}, nil
	}
	defer dbPool.Close()
	
	// 리포지토리 생성
	restaurantRepo := repository.NewRestaurantRepository(dbPool)
	
	// 서비스 생성
	s3Svc := service.NewPresignedURL(cfg, s3Client, presignClient)
	adminSvc := service.NewAdmin(cfg, restaurantRepo)
	
	// 핸들러 생성 및 요청 처리
	h := handler.NewHandler(cfg, s3Svc, adminSvc)
	return h.HandleRequest(ctx, request)
}

func main() {
	lambda.Start(handleRequest)
}
