package web


type CartItem struct {
	Produkid int `json:"produkid"`
	Quantity int `json:"quantity"`
}

type CartCheckoutPayload struct {
	Items []CartItem `json:"items" validate:"required"`
}

type JoinTableWeb struct {
	Id int `json:"id"`
	OrderId int `json:"orderId"`
	ProdukId int `json:"produkId"`
	UserId int `json:"userId"`
	Quantity int `json:"quantity"`
	TotalPrice float64 `json:"totalPrice"`
}