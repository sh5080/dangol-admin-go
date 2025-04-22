package routes

import (
	"context"
	"database/sql"
	appCtx "lambda-go/pkg/contexts"
	middleware "lambda-go/pkg/middlewares"
	"lambda-go/pkg/utils"
	"strings"

	config "lambda-go/pkg/configs"
	handler "lambda-go/pkg/handlers"
	adminHandler "lambda-go/pkg/handlers/admin"
	publicHandler "lambda-go/pkg/handlers/public"
	adminService "lambda-go/pkg/services/admin"
	publicService "lambda-go/pkg/services/public"

	"github.com/aws/aws-lambda-go/events"
)

const (
	NoAuth AuthType = iota
	DefaultAuth
	SessionAuth
)

type AuthType int

type Route struct {
	Path     string
	Method   string
	Handler  func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	AuthType AuthType
}

type Router interface {
	AddRoute(route Route)
	Handle(ctx context.Context, request events.APIGatewayProxyRequest, db *sql.DB) (events.APIGatewayProxyResponse, *utils.AppError)
}

type router struct {
	routes []Route
}

func NewRouter() Router {
	return &router{
		routes: []Route{},
	}
}

func (r *router) AddRoute(route Route) {
	r.routes = append(r.routes, route)
}

func (r *router) Handle(ctx context.Context, request events.APIGatewayProxyRequest, db *sql.DB) (events.APIGatewayProxyResponse, *utils.AppError) {
	// OPTIONS 메서드는 모든 경로에서 동일하게 처리
	if request.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{}, nil
	}

	// 라우트 일치 확인
	for _, route := range r.routes {
		pathMatches := false
		params := make(appCtx.Params)

		// 정확한 경로 매치 확인
		if route.Path == request.Path {
			pathMatches = true
		} else if strings.HasSuffix(route.Path, "*") {
			// 와일드카드 패턴 지원
			prefix := route.Path[:len(route.Path)-1]
			pathMatches = len(request.Path) >= len(prefix) && request.Path[:len(prefix)] == prefix
		} else if strings.Contains(route.Path, "{") && strings.Contains(route.Path, "}") {
			// 파라미터 패턴 지원 (예: /admin/restaurant/request/{id}/process)
			routeParts := strings.Split(route.Path, "/")
			requestParts := strings.Split(request.Path, "/")

			if len(routeParts) == len(requestParts) {
				matches := true
				for i, part := range routeParts {
					if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
						// 파라미터 부분은 항상 매치
						paramName := part[1 : len(part)-1]
						params[paramName] = requestParts[i]
					} else if part != requestParts[i] {
						matches = false
						break
					}
				}
				pathMatches = matches
			}
		}

		// 메서드 확인 (OPTIONS는 이미 처리했으므로 제외)
		methodMatches := route.Method == "" || route.Method == request.HTTPMethod

		if pathMatches && methodMatches {
			// 파라미터를 컨텍스트에 저장
			paramCtx := context.WithValue(ctx, appCtx.ParamsKey, params)

			// 인증 방식에 따라 처리
			switch route.AuthType {
			case SessionAuth:
				if db == nil {
					return events.APIGatewayProxyResponse{}, utils.InternalServerError("인증 처리를 위한 데이터베이스 연결이 없습니다", nil)
				}
				// 세션 인증 미들웨어 적용
				return middleware.SessionAuth(db, route.Handler)(paramCtx, request)

			case DefaultAuth:
				// 기본 JWT 인증 미들웨어 적용
				return middleware.DefaultAuth(route.Handler)(paramCtx, request)

			case NoAuth:
				// 인증 없음
				response, err := route.Handler(paramCtx, request)
				if err != nil {
					return events.APIGatewayProxyResponse{}, utils.InternalServerError("처리 중 오류가 발생했습니다", err)
				}
				return response, nil
			}
		}
	}

	// 일치하는 라우트가 없는 경우
	return events.APIGatewayProxyResponse{}, utils.NotFound("요청한 API를 찾을 수 없습니다")
}

// SetupRouter는 핸들러를 생성하고 라우트를 등록합니다.
func SetupRouter(
	ctx context.Context,
	cfg *config.Config,
	s3Svc *publicService.S3Service,
	adminSvc *adminService.RestaurantService,
	db *sql.DB,
) (Router, func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, *utils.AppError)) {
	// 기본 핸들러 생성
	h := handler.NewHandler(cfg, s3Svc, adminSvc)

	// 도메인별 핸들러 생성
	adminHandler := &adminHandler.AdminHandler{Handler: h}
	s3Handler := &publicHandler.S3Handler{Handler: h}

	// 라우터 생성
	router := NewRouter()

	// 라우트 등록
	RegisterAdminRoutes(router, adminHandler)
	RegisterPublicRoutes(router, s3Handler)

	// 핸들러 함수 반환
	handleFunc := func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, *utils.AppError) {
		return router.Handle(ctx, request, db)
	}

	return router, handleFunc
}
