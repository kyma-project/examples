package model

type Order struct {
	ID    string  `json:"orderId"`
	Total float64 `json:"total"`
}
