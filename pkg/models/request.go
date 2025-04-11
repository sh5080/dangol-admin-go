package models

// PresignedURLRequest는 클라이언트로부터 받는 요청 구조체입니다.
type PresignedURLRequest struct {
	Bucket      string `json:"bucket"`      // S3 버킷 이름
	Key         string `json:"key"`         // 파일 경로/이름
	Method      string `json:"method"`      // HTTP 메서드 (GET/PUT)
	Duration    int64  `json:"duration"`    // URL 유효 기간(초)
	ContentType string `json:"contentType"` // 파일 콘텐츠 타입 (PUT 요청시 유효)
	MaxFileSize int64  `json:"maxFileSize"` // 최대 파일 크기 (바이트 단위, PUT 요청시 유효)
}

// PresignedURLResponse는 생성된 presigned URL을 포함한 응답 구조체입니다.
type PresignedURLResponse struct {
	URL          string `json:"url"`          // 생성된 presigned URL
	ContentType  string `json:"contentType"`  // 콘텐츠 타입 (PUT 요청시만 반환)
	MaxFileSize  int64  `json:"maxFileSize"`  // 최대 파일 크기 (PUT 요청시만 반환)
	ExpiresAt    int64  `json:"expiresAt"`    // URL 만료 시간 (Unix 타임스탬프)
}

 