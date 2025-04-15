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

// ProcessRestaurantRequest는 매장 생성 요청 처리 페이로드입니다.
type ProcessRestaurantRequest struct {
	Status       RequestStatus `json:"status" validate:"required,oneof=APPROVED REJECTED"`
	RejectReason *string       `json:"rejectReason,omitempty" validate:"required_if=Status REJECTED,omitempty"`
} 
 