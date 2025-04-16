package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	config "lambda-go/pkg/configs"
	"lambda-go/pkg/models"
	"lambda-go/pkg/utils"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type contextKey string
const ClaimsKey contextKey = "claims"

type AdminSession struct {
	UserID string
	Token  string
	IP     string
}

type AdminUser struct {
	UserID string
	Role   string
}

func SessionAuth(db *sql.DB, next func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, *utils.AppError) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, *utils.AppError) {
		// OPTIONS 요청은 검증 없이 통과
		if request.HTTPMethod == "OPTIONS" {
			response, err := next(ctx, request)
			if err != nil {
				return events.APIGatewayProxyResponse{}, utils.InternalServerError("처리 중 오류가 발생했습니다", err)
			}
			return response, nil
		}

		cookieHeader, ok := request.Headers["Cookie"]
		if !ok || cookieHeader == "" {
			return events.APIGatewayProxyResponse{}, utils.Unauthorized("인증 정보가 필요합니다")
		}

		sessionToken := extractCookieValue(cookieHeader, "Admin-Session")
		if sessionToken == "" {
			return events.APIGatewayProxyResponse{}, utils.Unauthorized("유효한 어드민 세션이 필요합니다")
		}
		authHeader, ok := request.Headers["Authorization"]
		if !ok || authHeader == "" {
			return events.APIGatewayProxyResponse{}, utils.Unauthorized("인증 정보가 필요합니다")
		}
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.VerifyToken(accessToken, config.NewConfig())
	
		if err != nil {
			return events.APIGatewayProxyResponse{}, utils.Unauthorized(fmt.Sprintf("토큰 검증 실패: %s", err.Error()))
		}
		if claims.Role != models.ADMIN {
			return events.APIGatewayProxyResponse{}, utils.Unauthorized("관리자 권한이 없습니다")
		}
		// 세션 토큰 검증
		clientIP, err := getClientIP(request)
		if err != nil {
			return events.APIGatewayProxyResponse{}, utils.Unauthorized(fmt.Sprintf("정상적인 로그인이 아닙니다. 새로운 환경에서 다시 시도해주세요: %s", err.Error()))
		}
		_, err = validateAdminSession(ctx, db, claims.UserID, sessionToken, clientIP)
		if err != nil {
			return events.APIGatewayProxyResponse{}, utils.Unauthorized(fmt.Sprintf("세션 검증 실패: %s", err.Error()))
		}

		// 다음 핸들러 호출
		response, err := next(ctx, request)
		if err != nil {
			return events.APIGatewayProxyResponse{}, utils.InternalServerError("처리 중 오류가 발생했습니다", err)
		}
		return response, nil
	}
}

// validateAdminSession은 세션 토큰의 유효성을 검증합니다
func validateAdminSession(ctx context.Context, db *sql.DB, userID string, token string, clientIP string) (bool, error) {
	var session AdminSession

	// 세션 쿼리 - JOIN을 사용하여 사용자 정보도 함께 가져옴
	query := `
		SELECT s.userId, s.token, s.ip
		FROM AdminSession s
		WHERE s.userId = ?
	`
	err := db.QueryRowContext(ctx, query, userID).Scan(&session.UserID, &session.Token, &session.IP)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("세션을 찾을 수 없습니다")
		}
		return false, fmt.Errorf("데이터베이스 오류: %s", err.Error())
	}
	if session.Token != token {
		return false, errors.New("세션 토큰 불일치")
	}

	if session.IP != clientIP {	    
	    return false, errors.New("세션 IP 불일치")
	}
	return true, nil
}

// 쿠키 값 추출 유틸리티 함수
func extractCookieValue(cookieHeader, name string) string {
	cookies := parseCookies(cookieHeader)
	return cookies[name]
}

// 쿠키 헤더 파싱 함수
func parseCookies(cookieHeader string) map[string]string {
	cookies := make(map[string]string)
	
	parts := strings.Split(cookieHeader, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) != 2 {
			continue
		}
		
		key := keyValue[0]
		value := keyValue[1]
		cookies[key] = value
	}
	
	return cookies
}

func getClientIP(request events.APIGatewayProxyRequest) (string, error) {
	// API Gateway 프록시 요청에서 클라이언트 IP 추출
	if ip, ok := request.Headers["X-Forwarded-For"]; ok && ip != "" {
		// X-Forwarded-For는 콤마로 구분된 IP 목록일 수 있으므로 첫 번째 IP 사용
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0]), nil
	}
	
	// 헤더에 없다면 요청 컨텍스트에서 추출 시도
	if request.RequestContext.Identity.SourceIP != "" {
		return request.RequestContext.Identity.SourceIP, nil
	}
	
	return "", errors.New("클라이언트 IP 추출 실패")
}


