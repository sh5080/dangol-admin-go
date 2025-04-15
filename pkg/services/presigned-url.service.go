package service

import (
	"context"
	"fmt"
	"time"

	config "lambda-go/pkg/configs"
	"lambda-go/pkg/models"
	"lambda-go/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// PresignedURLService는 Presigned URL 생성과 관련된 서비스를 제공합니다.
type PresignedURL struct {
	config        *config.Config
	s3Client      *s3.Client
	presignClient *s3.PresignClient
}

// NewPresignedURLService는 새 PresignedURLService 인스턴스를 생성합니다.
func NewPresignedURL(cfg *config.Config, s3Client *s3.Client, presignClient *s3.PresignClient) *PresignedURL {
	return &PresignedURL{
		config:        cfg,
		s3Client:      s3Client,
		presignClient: presignClient,
	}
}

// ValidateAndPreprocessRequest는 요청을 검증하고 기본값을 설정합니다.
func (s *PresignedURL) ValidateAndPreprocessRequest(req *models.PresignedURLRequest) error {
	// 버킷 검증
	if req.Bucket == "" {
		req.Bucket = s.config.DefaultBucket
		if req.Bucket == "" {
			return fmt.Errorf("bucket이 필요합니다")
		}
	}

	// 키 검증
	if req.Key == "" {
		return fmt.Errorf("key가 필요합니다")
	}

	// Method 기본값 설정
	if req.Method == "" {
		req.Method = "GET"
	} else if req.Method != "GET" && req.Method != "PUT" {
		return fmt.Errorf("지원하지 않는 HTTP 메서드입니다. GET 또는 PUT만 지원합니다")
	}

	// Duration 기본값 설정
	if req.Duration <= 0 {
		req.Duration = 3600 // 1시간 기본값
	}

	// ContentType 기본값 설정 (PUT 요청시)
	if req.Method == "PUT" && req.ContentType == "" {
		// 파일 확장자로부터 Content-Type 유추
		if ext := utils.GetFileExtension(req.Key); ext != "" {
			req.ContentType = utils.GetMimeTypeFromExtension(ext)
		} else {
			req.ContentType = "application/octet-stream"
		}
	}

	// 최대 파일 크기 기본값 설정 (PUT 요청시)
	if req.Method == "PUT" && req.MaxFileSize <= 0 {
		req.MaxFileSize = s.config.DefaultMaxFileSize
	}

	return nil
}

// GeneratePresignedURL은 주어진 요청에 대한 Presigned URL을 생성합니다.
func (s *PresignedURL) GeneratePresignedURL(ctx context.Context, req *models.PresignedURLRequest) (*models.PresignedURLResponse, error) {
	var presignedURL string
	expiresAt := time.Now().Add(time.Duration(req.Duration) * time.Second).Unix()

	switch req.Method {
	case "GET":
		presignReq, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(req.Bucket),
			Key:    aws.String(req.Key),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(req.Duration) * time.Second
		})
		if err != nil {
			return nil, fmt.Errorf("presigned URL 생성 실패: %w", err)
		}
		presignedURL = presignReq.URL

	case "PUT":
		putObjectInput := &s3.PutObjectInput{
			Bucket: aws.String(req.Bucket),
			Key:    aws.String(req.Key),
		}

		// Content-Type 설정
		if req.ContentType != "" {
			putObjectInput.ContentType = aws.String(req.ContentType)
		}

		presignReq, err := s.presignClient.PresignPutObject(ctx, putObjectInput, func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(req.Duration) * time.Second
		})
		if err != nil {
			return nil, fmt.Errorf("presigned URL 생성 실패: %w", err)
		}
		presignedURL = presignReq.URL
	}

	// 응답 생성
	resp := &models.PresignedURLResponse{
		URL:       presignedURL,
		ExpiresAt: expiresAt,
	}

	// PUT 요청일 경우 추가 정보 제공
	if req.Method == "PUT" {
		resp.ContentType = req.ContentType
		resp.MaxFileSize = req.MaxFileSize
	}

	return resp, nil
}