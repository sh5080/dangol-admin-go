package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	config "lambda-go/pkg/configs"
	handler "lambda-go/pkg/handlers"
	middleware "lambda-go/pkg/middlewares"
	repository "lambda-go/pkg/repositories"
	service "lambda-go/pkg/services"
	"lambda-go/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// AppErrorToResponse는 AppError를 API Gateway 응답으로 변환합니다
func appErrorToResponse(err *utils.AppError) events.APIGatewayProxyResponse {
	response, _ := utils.Error(err.StatusCode, err.Message)
	
	// CORS 헤더 추가
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}
	response.Headers["Access-Control-Allow-Origin"] = "*"
	response.Headers["Access-Control-Allow-Headers"] = "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token"
	response.Headers["Access-Control-Allow-Methods"] = "GET,POST,OPTIONS"
	
	return response
}

// 어댑터 함수: 의존성 주입을 통해 핸들러 호출
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 설정 로드
	cfg := config.NewConfig()
	
	// S3 클라이언트 초기화
	s3Client, presignClient, err := config.NewS3Client(ctx, cfg)
	if err != nil {
		log.Printf("S3 클라이언트 초기화 실패: %v", err)
		appErr := utils.InternalServerError("서버 초기화 중 오류가 발생했습니다", err)
		return appErrorToResponse(appErr), nil
	}
	
	// 데이터베이스 초기화 (pgxpool 사용)
	dbCfg := cfg.NewDBConfig()
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.Database, dbCfg.SSLMode)
	
	dbPool, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		log.Printf("데이터베이스 연결 실패: %v", err)
		appErr := utils.InternalServerError("데이터베이스 연결 중 오류가 발생했습니다", err)
		return appErrorToResponse(appErr), nil
	}
	defer dbPool.Close()
	
	// SQL DB 연결 생성 (미들웨어용)
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("SQL 데이터베이스 연결 실패: %v", err)
		appErr := utils.InternalServerError("데이터베이스 연결 중 오류가 발생했습니다", err)
		return appErrorToResponse(appErr), nil
	}
	defer sqlDB.Close()
	
	// 리포지토리 생성
	restaurantRepo := repository.NewRestaurantRepository(dbPool)
	
	// 서비스 생성
	s3Svc := service.NewPresignedURL(cfg, s3Client, presignClient)
	adminSvc := service.NewAdmin(cfg, restaurantRepo)
	
	// 핸들러 생성
	h := handler.NewHandler(cfg, s3Svc, adminSvc)
	
	// 어드민 경로일 경우 인증 미들웨어 적용
	if strings.HasPrefix(request.Path, "/admin/") {
		// 인증 미들웨어 적용
		response, appErr := middleware.SessionAuth(sqlDB, h.HandleRequest)(ctx, request)
		if appErr != nil {
			return appErrorToResponse(appErr), nil
		}
		return response, nil
	}
	
	// 일반 경로는 직접 처리
	response, err := h.HandleRequest(ctx, request)
	if err != nil {
		// 일반 에러는 내부 서버 오류로 처리
		appErr := utils.InternalServerError("처리 중 오류가 발생했습니다", err)
		return appErrorToResponse(appErr), nil
	}
	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
