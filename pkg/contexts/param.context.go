package context

import (
	"context"
	"errors"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

type contextKey string

const ParamsKey contextKey = "urlParams"

// URL 파라미터 타입
type Params map[string]string

// GetParams는 컨텍스트에서 URL 파라미터를 가져옵니다
func GetParams(ctx context.Context) Params {
	params, ok := ctx.Value(ParamsKey).(Params)
	if !ok {
		return make(Params)
	}
	return params
}

// GetParam은 컨텍스트에서 특정 URL 파라미터를 가져옵니다
func GetParam(ctx context.Context, name string) string {
	return GetParams(ctx)[name]
}

var ErrInvalidParam = errors.New("잘못된 파라미터 형식")

// GetIntParam은 쿼리 파라미터에서 int 값을 추출합니다
func GetIntParam(request events.APIGatewayProxyRequest, paramName string, defaultValue int) int {
	if paramStr := request.QueryStringParameters[paramName]; paramStr != "" {
		if value, err := strconv.Atoi(paramStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// GetStringParam은 쿼리 파라미터에서 string 값을 추출합니다
func GetStringParam(request events.APIGatewayProxyRequest, paramName string, defaultValue string) string {
	if paramStr := request.QueryStringParameters[paramName]; paramStr != "" {
		return paramStr
	}
	return defaultValue
}

// GetBoolParam은 쿼리 파라미터에서 bool 값을 추출합니다
func GetBoolParam(request events.APIGatewayProxyRequest, paramName string, defaultValue bool) bool {
	if paramStr := request.QueryStringParameters[paramName]; paramStr != "" {
		if value, err := strconv.ParseBool(paramStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// ParsePaginationParams는 페이지네이션 관련 파라미터를 파싱합니다
func ParsePaginationParams(request events.APIGatewayProxyRequest) (page, pageSize int) {
	page = GetIntParam(request, "page", 1)
	if page < 1 {
		page = 1
	}

	pageSize = GetIntParam(request, "pageSize", 10)
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10 // 기본값으로 제한
	}

	return page, pageSize
}
