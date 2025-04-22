package dtos

import "lambda-go/pkg/models"

// RestaurantRequestQuery는 매장 요청 조회를 위한 쿼리 파라미터 DTO입니다.
type RestaurantRequestQuery struct {
	Page     int                             `json:"page"`
	PageSize int                             `json:"pageSize"`
	Status   *models.RestaurantRequestStatus `json:"status,omitempty"`
}
