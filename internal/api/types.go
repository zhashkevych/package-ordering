package api

type SetPacksRequest struct {
	PackSizes []int `json:"packSizes"`
}

type CalculateRequest struct {
	Amount int `json:"amount"`
}

type CalculateResponse struct {
	Amount     int         `json:"amount"`
	PackSizes  []int       `json:"packSizes"`
	Allocation map[int]int `json:"allocation"`
	TotalItems int         `json:"totalItems"`
	TotalPacks int         `json:"totalPacks"`
	Overfill   int         `json:"overfill"`
}
