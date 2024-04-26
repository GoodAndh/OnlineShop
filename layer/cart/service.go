package cart

import (
	"context"
	"database/sql"
	"ddd2/exception"
	"ddd2/layer/orders"
	"ddd2/model/web"
	"ddd2/utils"
)

type ServiceImpl struct {
	Orders orders.Service
	Db     *sql.DB
	repo   Repository
}
type Service interface {
	CreateOrders(ctx context.Context, ps []web.Produk, cart []web.CartItem, userid int) error
	GetOrdersList(ctx context.Context, userId int, orderId []int) ([]CartService, error)
	GetOrderId(ctx context.Context, userId int) ([]int, error)
	UpdateOrders(ctx context.Context, ord web.Orders) error
}

func NewService(orders orders.Service, db *sql.DB, repo Repository) Service {
	return &ServiceImpl{
		Orders: orders,
		Db:     db,
		repo:   repo,
	}
}

func checkIfinStock(pm map[int]web.Produk, cart []web.CartItem) error {
	if len(cart) <= 0 {
		return exception.ErrCartEmpty
	}

	for _, items := range cart {
		product, ok := pm[items.Produkid]
		if !ok {
			return exception.ErrIdNotFound
		}

		if product.Quantity < items.Quantity {
			return exception.ErrQuantityMoreThanStock
		}
	}
	return nil
}

// calculate total price
func calculateTotalPrice(pm map[int]web.Produk, cart []web.CartItem) float64 {
	var total float64
	for _, items := range cart {
		product := pm[items.Produkid]
		total += float64(product.Harga) * float64(items.Quantity)
	}
	return total
}

func (s *ServiceImpl) CreateOrders(ctx context.Context, ps []web.Produk, cart []web.CartItem, userid int) error {

	//store ps into map
	productMap := make(map[int]web.Produk)
	for _, product := range ps {
		productMap[product.Id] = product
	}

	//calculate total price
	totalPrice := calculateTotalPrice(productMap, cart)

	//check if Cart is in stock
	if err := checkIfinStock(productMap, cart); err != nil {
		return err
	}

	orderId, err := s.Orders.CreateOrders(ctx, web.Orders{
		Id:      0,
		Userid:  userid,
		Status:  "pending",
		Address: "some address",
	})
	if err != nil {
		return err
	}
	for _, items := range cart {
		err = s.Orders.CreateOrdersItems(ctx, web.OrdersItems{
			Userid:   userid,
			Produkid: items.Produkid,
			Orderid:  orderId,
			Quantity: items.Quantity,
			Price:    totalPrice,
		})
		if err != nil {
			return err
		}

	}

	return nil

}

func (s *ServiceImpl) GetOrdersList(ctx context.Context, userId int, orderId []int) ([]CartService, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	cart, err := s.repo.GetOrdersList(ctx, tx, userId, orderId)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *ServiceImpl) GetOrderId(ctx context.Context, userId int) ([]int, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	ordId, err := s.repo.GetOrderId(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return ordId, nil
}

func (s *ServiceImpl) UpdateOrders(ctx context.Context, ord web.Orders) error {
	err:=s.Orders.UpdateOrders(ctx,ord)
	return err
}
