package routes

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// RestaurantHandler는 Restaurant 관련 핸들러 인터페이스
type RestaurantHandler interface {
	GetRestaurantRequests(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	ProcessRestaurantRequest(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func RegisterAdminRoutes(router Router, h RestaurantHandler) {
	// 매장 생성 요청 목록 조회 API
	router.AddRoute(Route{
		Path:     "/admin/restaurant/request",
		Method:   "GET",
		Handler:  h.GetRestaurantRequests,
		AuthType: NoAuth,
	})

	// 매장 생성 요청 처리 API
	router.AddRoute(Route{
		Path:     "/admin/restaurant/request/{id}/process",
		Method:   "POST",
		Handler:  h.ProcessRestaurantRequest,
		AuthType: NoAuth,
	})
}
