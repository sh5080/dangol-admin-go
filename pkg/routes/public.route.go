package routes

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

// S3Handler는 S3 관련 핸들러 인터페이스
type S3Handler interface {
	GetPresignedURL(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func RegisterPublicRoutes(router Router, h S3Handler) {
	// 업로드용 Presigned URL 발급 API
	router.AddRoute(Route{
		Path:     "/s3/presigned-url",
		Method:   "GET",
		Handler:  h.GetPresignedURL,
		AuthType: NoAuth,
	})
}
