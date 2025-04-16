package main

import (
	"context"
	"database/sql"
	"log"

	config "lambda-go/pkg/configs"
	handler "lambda-go/pkg/handlers"
	repository "lambda-go/pkg/repositories"
	"lambda-go/pkg/routes"
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
	connStr := dbCfg.DatabaseURL
	
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
	
	restaurantRepo := repository.NewRestaurantRepository(dbPool)

	s3Svc := service.NewPresignedURL(cfg, s3Client, presignClient)
	adminSvc := service.NewAdmin(cfg, restaurantRepo)
	
	h := handler.NewHandler(cfg, s3Svc, adminSvc)
	
	router := routes.NewRouter()
	
	routes.RegisterAdminRoutes(router, h)
	routes.RegisterPublicRoutes(router, h)
	
	// 라우터를 통한 요청 처리
	response, appErr := router.Handle(ctx, request, sqlDB)
	if appErr != nil {
		return appErrorToResponse(appErr), nil
	}
	
	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
