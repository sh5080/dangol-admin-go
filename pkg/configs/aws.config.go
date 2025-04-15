package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewS3Client는 AWS S3 클라이언트와 Presign 클라이언트를 생성합니다.
func NewS3Client(ctx context.Context, cfg *Config) (*s3.Client, *s3.PresignClient, error) {
	// 기본 옵션 슬라이스 생성
	options := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.AWSRegion),
	}

	// 명시적인 AWS 자격 증명이 제공된 경우 사용
	if cfg.AWSAccessKeyID != "" && cfg.AWSSecretAccessKey != "" {
		options = append(options, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKeyID,
			cfg.AWSSecretAccessKey,
			"", // Session token (선택적)
		)))
	}

	// AWS 설정 로드
	awsCfg, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return nil, nil, fmt.Errorf("AWS 설정 로드 실패: %w", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)
	presignClient := s3.NewPresignClient(s3Client)
	
	return s3Client, presignClient, nil
}