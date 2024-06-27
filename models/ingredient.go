package models

type Ingredient struct {
	Name     string  `json:"name"`
	Unit     string  `json:"unit"`
	Quantity float64 `json:"quantity"`
}
