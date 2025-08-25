package models

type Order struct {
	ID     string   `json:"id"`
	Amount float64  `json:"amount"`
	Items  []string `json:"items"`
}
