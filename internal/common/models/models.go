package models

type Metric struct {
	ID    string   `json:"id"`
	Mtype string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
