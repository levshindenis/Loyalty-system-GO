package models

import "time"

type Order struct {
	OrderId    string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	WithDrawn float64 `json:"withdrawn"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Withdraw struct {
	OrderId string  `json:"order"`
	Summ    float64 `json:"sum"`
}

type OutPoints struct {
	OrderId     string    `json:"order"`
	Summ        float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

// Для клиента

type Product struct {
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type RegOrder struct {
	OrderId  string    `json:"order"`
	Products []Product `json:"goods"`
}

type Good struct {
	Match      string  `json:"match"`
	Reward     float64 `json:"reward"`
	RewardType string  `json:"reward_type"`
}
