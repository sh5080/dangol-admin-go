package models

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type APIErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// RestaurantRequestsResponse는 매장 생성 요청 목록 응답입니다.
type RestaurantRequestsResponse struct {
	Requests []RestaurantRequest `json:"requests"`
	Total    int                 `json:"total"`
}

// PresignedURLResponse는 생성된 presigned URL을 포함한 응답 구조체입니다.
type PresignedURLResponse struct {
	URL          string `json:"url"`          // 생성된 presigned URL
	ContentType  string `json:"contentType"`  // 콘텐츠 타입 (PUT 요청시만 반환)
	MaxFileSize  int64  `json:"maxFileSize"`  // 최대 파일 크기 (PUT 요청시만 반환)
	ExpiresAt    int64  `json:"expiresAt"`    // URL 만료 시간 (Unix 타임스탬프)
}