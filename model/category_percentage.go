package model

//CategoryPercentage summarises a category
type CategoryPercentage struct {
	CategoryID int64   `json:"category_id"`
	Available  int64   `json:"active"`
	Percentage float64 `json:"percentage"`
}
