package structs

import "time"

type Order struct {
	OrderId    string  `json:"order_id"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type Balance struct {
	Current   float32 `json:"current"`
	WithDrawn float32 `json:"withdrawn"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Withdraw struct {
	OrderId string  `json:"order"`
	Summ    float32 `json:"sum"`
}

type OutPoints struct {
	OrderId     string    `json:"order_id"`
	Summ        float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
