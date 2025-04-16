package utils

import (
	"fmt"
	config "lambda-go/pkg/configs"
	"lambda-go/pkg/models"

	"github.com/golang-jwt/jwt/v4"
)

// Claims는 JWT 토큰의 페이로드 구조를 정의합니다
type Claims struct {
	UserID string      `json:"userId"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// VerifyToken JWT 토큰을 검증하고 클레임을 반환합니다
func VerifyToken(tokenString string, cfg *config.Config) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, 
		&Claims{}, 
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, Unauthorized(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
			}
			
			if token.Method.Alg() != "HS256" {
				return nil, Unauthorized("algorithm must be HS256")
			}
			
			return []byte(cfg.JWTSecret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	// 에러가 있는 경우 처리
	if err != nil {
		// 토큰 만료 에러 확인
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// 토큰이 만료된 경우 정상적으로 claims 추출 (이후 세션 검증)
				if claims, ok := token.Claims.(*Claims); ok {
					return claims, nil
				}
			}
		}
		return nil, Unauthorized(fmt.Sprintf("토큰 검증 실패: %s", err.Error()))
	}
	
	// 토큰이 유효한 경우
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, Unauthorized("유효하지 않은 토큰입니다")
	}
	
	return claims, nil
}
