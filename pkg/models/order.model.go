package models

import "time"

// OrderStatus는 주문 상태를 나타내는 열거형입니다.
type OrderStatus string

const (
	ORDER_NEW        OrderStatus = "NEW"
	ORDER_PROCESSING OrderStatus = "PROCESSING"
	ORDER_COMPLETED  OrderStatus = "COMPLETED"
	ORDER_REJECTED   OrderStatus = "REJECTED"
	ORDER_CANCELLED  OrderStatus = "CANCELLED"
)

// DeliveryStatus는 배달 상태를 나타내는 열거형입니다.
type DeliveryStatus string

const (
	NOT_CALLED         DeliveryStatus = "NOT_CALLED" // 배달 준비 중
	DELIVERING         DeliveryStatus = "DELIVERING" // 배달 중
	DELIVERY_COMPLETED DeliveryStatus = "COMPLETED"  // 배달 완료
)

// DeliveryType은 배달 유형을 나타내는 열거형입니다.
type DeliveryType string

const (
	SELF     DeliveryType = "SELF"     // 자체배달
	DELIVERY DeliveryType = "DELIVERY" // 배달대행
	PICKUP   DeliveryType = "PICKUP"   // 픽업
)

// Order는 주문 모델입니다.
type Order struct {
	ID           string         `json:"id" db:"id"`
	RestaurantID string         `json:"restaurantId" db:"restaurantId"`
	CustomerID   string         `json:"customerId" db:"customerId"`
	Status       OrderStatus    `json:"status" db:"status"`
	CreatedAt    time.Time      `json:"createdAt" db:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt" db:"updatedAt"`
	Restaurant   Restaurant     `json:"restaurant,omitempty"`
	Menus        []OrderMenu    `json:"menus,omitempty"`
	Customer     User           `json:"customer,omitempty"`
	Delivery     *OrderDelivery `json:"delivery,omitempty"`
}

// OrderDelivery는 주문 배달 정보 모델입니다.
type OrderDelivery struct {
	ID            string         `json:"id" db:"id"`
	OrderID       string         `json:"orderId" db:"orderId"`
	Status        DeliveryStatus `json:"status" db:"status"`
	Type          *DeliveryType  `json:"type,omitempty" db:"type"`
	DeliveryInfo  interface{}    `json:"deliveryInfo,omitempty" db:"deliveryInfo"` // JSON 데이터
	EstimatedTime *int           `json:"estimatedTime,omitempty" db:"estimatedTime"`
	CreatedAt     time.Time      `json:"createdAt" db:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt" db:"updatedAt"`
	Order         Order          `json:"order,omitempty"`
}

// OrderMenu는 주문 메뉴 모델입니다.
type OrderMenu struct {
	OrderID  string         `json:"orderId" db:"orderId"`
	MenuID   string         `json:"menuId" db:"menuId"`
	Quantity int            `json:"quantity" db:"quantity"`
	Order    Order          `json:"order,omitempty"`
	Menu     RestaurantMenu `json:"menu,omitempty"`
}
