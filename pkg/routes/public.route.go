package routes

import (
	handler "lambda-go/pkg/handlers"
)

func RegisterPublicRoutes(router Router, h *handler.Handler) {
	router.AddRoute(Route{
		Path:     "/presigned-url",
		Method:   "POST",
		Handler:  h.PresignedURLHandleRequest,
		AuthType: DefaultAuth,  
	})
}
