package routes

import (
	handler "lambda-go/pkg/handlers"
)

func RegisterAdminRoutes(router Router, h *handler.Handler) {
	// 매장 생성 요청 목록 조회 API
	router.AddRoute(Route{
		Path:     "/admin/restaurant/request",
		Method:   "GET",
		Handler:  h.GetRestaurantRequests,
		AuthType: SessionAuth,
	})

	// 매장 생성 요청 처리 API
	router.AddRoute(Route{
		Path:     "/admin/restaurant/request/{id}/process",
		Method:   "POST",
		Handler:  h.ProcessRestaurantRequest,
		AuthType: SessionAuth,
	})
}


