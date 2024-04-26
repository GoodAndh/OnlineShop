package cart

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type CartService struct {
	Id         int     `json:"id"`
	Orderid    int     `json:"order_id"`
	Produkid   int     `json:"produk_id"`
	NameProduk string  `json:"produk_name"`
	Quantity   int     `json:"jumlah_pesanan"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
	SemuaHarga float64 `json:"semua_harga"`
}

type Repository interface {
	GetOrdersList(ctx context.Context, tx *sql.Tx, userId int, orderId []int) ([]CartService, error)
	GetOrderId(ctx context.Context, tx *sql.Tx, userId int) ([]int, error)
}
type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}
func (r *RepositoryImpl) GetOrdersList(ctx context.Context, tx *sql.Tx, userId int, orderId []int) ([]CartService, error) {

	placeholders := strings.Repeat("?,", len(orderId))

	query := fmt.Sprintf("select orders_items.id as id ,orders_items.orderid as order_id,orders_items.produkid as produk_id,produk.produkname as name_produk ,orders_items.quantity as jumlah_pesanan,produk.harga as harga_satuan, orders_items.price as total_harga from orders_items join produk on orders_items.produkid=produk.id where  orders_items.userid = %v and orders_items.orderid in (%s%v) ", userId, placeholders, 0)

	//convert slice of orderid into []interface
	args := make([]interface{}, len(orderId))
	for i, v := range orderId {
		args[i] = v
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	cart := make([]CartService, len(orderId))
	defer rows.Close()
	for rows.Next() {
		crt := &CartService{}
		err := rows.Scan(&crt.Id, &crt.Orderid, &crt.Produkid, &crt.NameProduk, &crt.Quantity, &crt.Price, &crt.TotalPrice)
		if err != nil {
			return nil, err
		}
		cart = append(cart, *crt)
	}
	rows,err=tx.QueryContext(ctx,fmt.Sprintf("select sum(price) from orders_items where orderid in (%s%v)",placeholders,0),args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		crt := &CartService{}
		err := rows.Scan(&crt.SemuaHarga)
		if err != nil {
			return nil, err
		}
		cart = append(cart, *crt)
	}

	

	return cart, nil
}

func (r *RepositoryImpl) GetOrderId(ctx context.Context, tx *sql.Tx, userId int) ([]int, error) {
	rows, err := tx.QueryContext(ctx, "select orderid from orders_items where userid = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orderid := make([]int, 0)

	for rows.Next() {
		var ordid int
		err := rows.Scan(&ordid)
		if err != nil {
			return nil, err
		}
		orderid = append(orderid, ordid)
	}

	return orderid, nil

}
