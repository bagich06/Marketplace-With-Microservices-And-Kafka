package models

type Order struct {
	ID          int     `json:"id"`
	ProductName string  `json:"product_name"`
	ProductID   int     `json:"product_id"`
	SupplierID  int     `json:"supplier_id"`
	ClientID    int     `json:"client_id"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}
