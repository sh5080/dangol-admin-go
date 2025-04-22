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

// RestaurantRequestStatus는 매장 요청 상태를 나타내는 열거형입니다.
type RestaurantRequestStatus string

const (
	PENDING  RestaurantRequestStatus = "PENDING"
	APPROVED RestaurantRequestStatus = "APPROVED"
	REJECTED RestaurantRequestStatus = "REJECTED"
)

// RestaurantStatus는 매장 상태를 나타내는 열거형입니다.
type RestaurantStatus string

const (
	REQUESTED RestaurantStatus = "REQUESTED"
	HIDDEN    RestaurantStatus = "HIDDEN"
	OPEN      RestaurantStatus = "OPEN"
	CLOSED    RestaurantStatus = "CLOSED"
)

// RestaurantRequestType는 매장 요청 유형을 나타내는 열거형입니다.
type RestaurantRequestType string

const (
	CREATE RestaurantRequestType = "CREATE"
	UPDATE RestaurantRequestType = "UPDATE"
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
	ID                 string              `json:"id" db:"id"`
	Name               string              `json:"name" db:"name"`
	Description        *string             `json:"description,omitempty" db:"description"`
	Address            string              `json:"address" db:"address"`
	PhoneNumber        string              `json:"phoneNumber" db:"phoneNumber"`
	OwnerID            string              `json:"ownerId" db:"ownerId"`
	AddressDescription *string             `json:"addressDescription,omitempty" db:"addressDescription"`
	EventDescription   *string             `json:"eventDescription,omitempty" db:"eventDescription"`
	Holiday            *string             `json:"holiday,omitempty" db:"holiday"`
	ParkingAvailable   bool                `json:"parkingAvailable" db:"parkingAvailable"`
	ParkingDescription *string             `json:"parkingDescription,omitempty" db:"parkingDescription"`
	DeliveryAvailable  bool                `json:"deliveryAvailable" db:"deliveryAvailable"`
	Status             RestaurantStatus    `json:"status" db:"status"`
	CreatedAt          time.Time           `json:"createdAt" db:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt" db:"updatedAt"`
	DeletedAt          *time.Time          `json:"deletedAt,omitempty" db:"deletedAt"`
	BusinessHours      []BusinessHour      `json:"businessHours,omitempty"`
	Tags               []RestaurantTag     `json:"tags,omitempty"`
	Requests           []RestaurantRequest `json:"requests,omitempty"`
	Menus              []RestaurantMenu    `json:"menus,omitempty"`
	Images             []RestaurantImage   `json:"images,omitempty"`
	Business           *RestaurantBusiness `json:"business,omitempty"`
	Orders             []Order             `json:"orders,omitempty"`
	Owner              *User               `json:"owner,omitempty"`
}

// RestaurantBusiness 매장 사업자 정보 모델입니다.
type RestaurantBusiness struct {
	RestaurantID    string     `json:"restaurantId" db:"restaurantId"`
	Name            string     `json:"name" db:"name"`
	LicenseImageUrl string     `json:"licenseImageUrl" db:"licenseImageUrl"`
	LicenseNumber   string     `json:"licenseNumber" db:"licenseNumber"`
	Restaurant      Restaurant `json:"restaurant,omitempty"`
}

// RestaurantImage는 매장 이미지 모델입니다.
type RestaurantImage struct {
	ID           string     `json:"id" db:"id"`
	RestaurantID string     `json:"restaurantId" db:"restaurantId"`
	ImageUrl     string     `json:"imageUrl" db:"imageUrl"`
	Restaurant   Restaurant `json:"restaurant,omitempty"`
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
	OrderMenus   []OrderMenu `json:"orderMenus,omitempty"`
}

// BusinessHour는 영업 시간 모델입니다.
type BusinessHour struct {
	ID           int        `json:"id" db:"id"`
	RestaurantID string     `json:"restaurantId" db:"restaurantId"`
	OpenTime     string     `json:"openTime" db:"openTime"`   // HHmm 형식
	CloseTime    string     `json:"closeTime" db:"closeTime"` // HHmm 형식
	DayOfWeek    DayOfWeek  `json:"dayOfWeek" db:"dayOfWeek"`
	Restaurant   Restaurant `json:"restaurant,omitempty"`
}

// RestaurantRequest는 매장 생성/수정 요청 모델입니다.
type RestaurantRequest struct {
	ID                      int                     `json:"id" db:"id"`
	RestaurantID            string                  `json:"restaurantId" db:"restaurantId"`
	UserID                  string                  `json:"userId" db:"userId"`
	Name                    *string                 `json:"name,omitempty" db:"name"`
	BusinessLicenseImageUrl *string                 `json:"businessLicenseImageUrl,omitempty" db:"businessLicenseImageUrl"`
	BusinessLicenseNumber   *string                 `json:"businessLicenseNumber,omitempty" db:"businessLicenseNumber"`
	RejectReason            *string                 `json:"rejectReason,omitempty" db:"rejectReason"`
	Type                    RestaurantRequestType   `json:"type" db:"type"`
	Status                  RestaurantRequestStatus `json:"status" db:"status"`
	CreatedAt               time.Time               `json:"createdAt" db:"createdAt"`
	UpdatedAt               time.Time               `json:"updatedAt" db:"updatedAt"`
	DeletedAt               *time.Time              `json:"deletedAt,omitempty" db:"deletedAt"`
	User                    *User                   `json:"user,omitempty"`
	Restaurant              *Restaurant             `json:"restaurant,omitempty"`
}
