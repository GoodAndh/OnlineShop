package domain


type Orders struct {
	Id int
	Userid int
	Status string
	Address string
}


type OrdersItems struct {
	Id int
	Userid int
	Produkid int
	Orderid int
	Quantity int
	Price float64
}