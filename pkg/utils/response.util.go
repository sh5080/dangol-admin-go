package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"lambda-go/pkg/models"

	"github.com/aws/aws-lambda-go/events"
)

func Success(statusCode int, data interface{}) (events.APIGatewayProxyResponse, error) {
	response := models.APIResponse{
		Status: "success",
		Data:   data,
	}
	
	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

func Error(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	errorResponse := models.APIErrorResponse{
		Status:  "error",
		Message: message,
	}
	
	body, err := json.Marshal(errorResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil
}

// AppError는 애플리케이션에서 발생하는 모든 에러를 표현합니다.
type AppError struct {
    StatusCode int    // HTTP 상태 코드
    Code       string // 에러 코드 (선택적)
    Message    string // 에러 메시지
    Err        error  // 원본 에러 (선택적)
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

// Unwrap은 원본 에러를 반환합니다.
func (e *AppError) Unwrap() error {
    return e.Err
}

// 자주 사용되는 에러 생성 함수
func BadRequest(message string, err ...error) *AppError {
    var original error
    if len(err) > 0 {
        original = err[0]
    }
    return &AppError{
        StatusCode: http.StatusBadRequest,
        Message:    message,
        Err:        original,
    }
}

func Unauthorized(message string, err ...error) *AppError {
    var original error
    if len(err) > 0 {
        original = err[0]
    }
    return &AppError{
        StatusCode: http.StatusUnauthorized,
        Message:    message,
        Err:        original,
    }
}

func Forbidden(message string, err ...error) *AppError {
    var original error
    if len(err) > 0 {
        original = err[0]
    }
    return &AppError{
        StatusCode: http.StatusForbidden,
        Message:    message,
        Err:        original,
    }
}

func NotFound(message string, err ...error) *AppError {
    var original error
    if len(err) > 0 {
        original = err[0]
    }
    return &AppError{
        StatusCode: http.StatusNotFound,
        Message:    message,
        Err:        original,
    }
}

func InternalServerError(message string, err ...error) *AppError {
    var original error
    if len(err) > 0 {
        original = err[0]
    }
    return &AppError{
        StatusCode: http.StatusInternalServerError,
        Message:    message,
        Err:        original,
    }
}