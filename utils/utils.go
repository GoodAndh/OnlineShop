package utils

import (
	"database/sql"
	"ddd2/model/domain"
	"ddd2/model/web"
)

func CommitOrRollback(tx *sql.Tx) error {
	err := recover()
	if err != nil {
		return tx.Rollback()
	} else {

		return tx.Commit()
	}
}

func ConvertUserToWeb(user *domain.User) *web.User {
	return &web.User{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Name:     user.Name,
	}
}

func ConvertProdukToWeb(produk *domain.Produk) *web.Produk {
	return &web.Produk{
		Id:         produk.Id,
		ProdukName: produk.ProdukName,
		Deskripsi:  produk.Deskripsi,
		Category:   produk.Category,
		UserId:     produk.UserId,
		Harga:      produk.Harga,
		Quantity:   produk.Quantity,
	}
}

func ConvertProdukSlice(produk []domain.Produk) []web.Produk {
	p := []web.Produk{}
	for _, v := range produk {
		p = append(p, *ConvertProdukToWeb(&v))
	}
	return p
}

func ConvertOrdersToWeb(ord *domain.Orders) *web.Orders {
	return &web.Orders{
		Id:      ord.Id,
		Userid:  ord.Userid,
		Status:  ord.Status,
		Address: ord.Address,
	}
}

func ConvertOrdersToSlice(ord []domain.Orders) []web.Orders {
	ORD := []web.Orders{}
	for _, v := range ord {
		ord = append(ord, domain.Orders(*ConvertOrdersToWeb(&v)))
	}
	return ORD
}

func ConvertItemsToWeb(ord *domain.OrdersItems) web.OrdersItems {
	return web.OrdersItems{
		Id:       ord.Id,
		Userid:   ord.Userid,
		Produkid: ord.Produkid,
		Orderid:  ord.Orderid,
		Quantity: ord.Quantity,
		Price:    ord.Price,
	}
}

func ConvertItemsToSlice(ord []domain.OrdersItems) []web.OrdersItems {
	ORD := []web.OrdersItems{}
	for _, v := range ord {
		ORD = append(ORD,ConvertItemsToWeb(&v))
	}
	return ORD
}
