package models

import (
	"time"
)

// DayOfWeek은 요일을 나타내는 열거형입니다.
type DayOfWeek string

const (
	MONDAY    DayOfWeek = "MON"
	TUESDAY   DayOfWeek = "TUE"
	WEDNESDAY DayOfWeek = "WED"
	THURSDAY  DayOfWeek = "THU"
	FRIDAY    DayOfWeek = "FRI"
	SATURDAY  DayOfWeek = "SAT"
	SUNDAY    DayOfWeek = "SUN"
)

// RequestStatus는 요청 상태를 나타내는 열거형입니다.
type RequestStatus string

const (
	PENDING  RequestStatus = "PENDING"
	APPROVED RequestStatus = "APPROVED"
	REJECTED RequestStatus = "REJECTED"
)

// RestaurantStatus는 매장 상태를 나타내는 열거형입니다.
type RestaurantStatus string

const (
	Requested RestaurantStatus = "REQUESTED"
	Hidden    RestaurantStatus = "HIDDEN"
	Open      RestaurantStatus = "OPEN"
	Closed    RestaurantStatus = "CLOSED"
)

// Tag는 태그 모델입니다.
type Tag struct {
	ID          int             `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	Description *string         `json:"description,omitempty" db:"description"`
	Restaurants []RestaurantTag `json:"restaurants,omitempty"`
}

// RestaurantTag는 매장과 태그의 다대다 관계를 나타내는 모델입니다.
type RestaurantTag struct {
	ID           int        `json:"id" db:"id"`
	RestaurantID string     `json:"restaurantId" db:"restaurantId"`
	TagID        int        `json:"tagId" db:"tagId"`
	Restaurant   Restaurant `json:"restaurant,omitempty"`
	Tag          Tag        `json:"tag,omitempty"`
}

// Restaurant는 매장 모델입니다.
type Restaurant struct {
	ID                      string              `json:"id" db:"id"`
	Name                    string              `json:"name" db:"name"`
	Description             *string             `json:"description,omitempty" db:"description"`
	BusinessLicenseImageURL string              `json:"businessLicenseImageUrl" db:"businessLicenseImageUrl"`
	BusinessLicenseNumber   string              `json:"businessLicenseNumber" db:"businessLicenseNumber"`
	Address                 string              `json:"address" db:"address"`
	PhoneNumber             string              `json:"phoneNumber" db:"phoneNumber"`
	OwnerID                 string              `json:"ownerId" db:"ownerId"`
	DeliveryAvailable       bool                `json:"deliveryAvailable" db:"deliveryAvailable"`
	Status                  RestaurantStatus    `json:"status" db:"status"`
	CreatedAt               time.Time           `json:"createdAt" db:"createdAt"`
	UpdatedAt               time.Time           `json:"updatedAt" db:"updatedAt"`
	DeletedAt               *time.Time          `json:"deletedAt,omitempty" db:"deletedAt"`
	BusinessHours           []BusinessHour      `json:"businessHours,omitempty"`
	Tags                    []RestaurantTag     `json:"tags,omitempty"`
	Requests                []RestaurantRequest `json:"requests,omitempty"`
	Menus                   []RestaurantMenu    `json:"menus,omitempty"`
}

// RestaurantMenu는 매장 메뉴 모델입니다.
type RestaurantMenu struct {
	ID           string      `json:"id" db:"id"`
	RestaurantID string      `json:"restaurantId" db:"restaurantId"`
	Name         string      `json:"name" db:"name"`
	Price        int         `json:"price" db:"price"`
	Description  *string     `json:"description,omitempty" db:"description"`
	ImageURL     *string     `json:"imageUrl,omitempty" db:"imageUrl"`
	Restaurant   Restaurant  `json:"restaurant,omitempty"`
}

// BusinessHour는 영업 시간 모델입니다.
type BusinessHour struct {
	ID           int        `json:"id" db:"id"`
	RestaurantID string     `json:"restaurantId" db:"restaurantId"`
	OpenTime     string     `json:"openTime" db:"openTime"` // HHmm 형식
	CloseTime    string     `json:"closeTime" db:"closeTime"` // HHmm 형식
	DayOfWeek    DayOfWeek  `json:"dayOfWeek" db:"dayOfWeek"`
	Restaurant   Restaurant `json:"restaurant,omitempty"`
}

// RestaurantRequest는 매장 생성 요청 모델입니다.
type RestaurantRequest struct {
	ID           int           `json:"id" db:"id"`
	RestaurantID string        `json:"restaurantId" db:"restaurantId"`
	UserID       string        `json:"userId" db:"userId"`
	RejectReason *string       `json:"rejectReason,omitempty" db:"rejectReason"`
	CreatedAt    time.Time     `json:"createdAt" db:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" db:"updatedAt"`
	DeletedAt    *time.Time    `json:"deletedAt,omitempty" db:"deletedAt"`
	Status       RequestStatus `json:"status" db:"status"`
	User         *User         `json:"user,omitempty"`
}

