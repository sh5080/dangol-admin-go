package context

import (
	"context"
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
