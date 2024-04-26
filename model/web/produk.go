package web

type Produk struct {
	Id         int    `json:"id"`
	ProdukName string `json:"name_produk"`
	Deskripsi  string `json:"deskripsi"`
	Category   string `json:"category"`
	UserId     int    `json:"user_id"`
	Harga      int    `json:"harga"`
	Quantity   int    `json:"quantity"`
}

type ProdukCreatePayload struct {
	ProdukName string `json:"name_produk" validate:"required"`
	Deskripsi  string `json:"deskripsi" validate:"required"`
	Category   string `json:"category" validate:"required,oneof=electric consumable etc"`
	UserId     int    `json:"user_id" validate:"required,number"`
	Harga      int    `json:"harga" validate:"required,number,gt=100"`
	Quantity   int    `json:"quantity" validate:"required,number,gt=0"`
}
