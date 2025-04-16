package config

import (
	"os"
	"strconv"
)

// Config는 애플리케이션 설정 값을 관리하는 구조체입니다.
type Config struct {
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	DefaultBucket      string
	DefaultMaxFileSize int64
	Environment        string
	JWTSecret          string
}

// 환경 변수로부터 기본값을 가져오는 함수
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// NewConfig는 환경 변수에서 설정을 로드하여 Config 구조체를 반환합니다.
func NewConfig() *Config {
	defaultMaxFileSizeStr := GetEnvOrDefault("DEFAULT_MAX_FILE_SIZE", "10485760")
	defaultMaxFileSize, err := strconv.ParseInt(defaultMaxFileSizeStr, 10, 64)
	if err != nil {
		defaultMaxFileSize = 10485760 // 10MB 기본값
	}

	return &Config{
		AWSRegion:          GetEnvOrDefault("AWS_REGION", "ap-northeast-2"),
		AWSAccessKeyID:     GetEnvOrDefault("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: GetEnvOrDefault("AWS_SECRET_ACCESS_KEY", ""),
		DefaultBucket:      GetEnvOrDefault("DEFAULT_BUCKET", ""),
		DefaultMaxFileSize: defaultMaxFileSize,
		Environment:        GetEnvOrDefault("ENV", "dev"),
		JWTSecret:          GetEnvOrDefault("JWT_SECRET", "1234567890abcdef"),
	}
}



