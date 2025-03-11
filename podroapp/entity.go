package podroapp

import "time"

type OrderStatus string

const (
	PendingOrderStatus      OrderStatus = "Pending"
	ProviderSeenOrderStatus OrderStatus = "ProviderSeen"
	PickedUpOrderStatus     OrderStatus = "PickedUp"
	InprogressOrderStatus   OrderStatus = "Inprogress"
	DeliveredOrderStatus    OrderStatus = "Delivered"
)

var OrderStatusStrings = map[OrderStatus]string{
	PendingOrderStatus:      "Pending",
	ProviderSeenOrderStatus: "ProviderSeen",
	PickedUpOrderStatus:     "PickedUp",
	InprogressOrderStatus:   "Inprogress",
	DeliveredOrderStatus:    "Delivered",
}

func (os OrderStatus) IsValid() bool {
	_, ok := OrderStatusStrings[os]
	return ok
}

type Customer struct {
	ID        uint
	Name      string
	Phone     string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Provider struct {
	ID        uint
	Name      string
	API       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID               uint
	ProviderID       uint
	CustomerID       uint
	CustomerName     string
	CustomerPhone    string
	CustomerAddress  string
	RecipientPhone   string
	RecipientName    string
	RecipientAddress string
	Status           OrderStatus
	CreatedAt        time.Time
	PickedUpAt       time.Time
	DeliveredAt      time.Time
	UpdatedAt        time.Time
}
