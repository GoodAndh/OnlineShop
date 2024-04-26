package web

type Orders struct {
	Id      int    `json:"id"`
	Userid  int    `json:"userid"`
	Status  string `json:"status"`
	Address string `json:"address"`
}

type OrdersItems struct {
	Id       int     `json:"id"`
	Userid   int     `json:"userid"`
	Produkid int     `json:"produkid"`
	Orderid  int     `json:"orderid"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
