package orders

import (
	"context"
	"database/sql"
	"ddd2/model/domain"
	"ddd2/model/web"
	"ddd2/utils"
)

type Service interface {
	GetOrders(ctx context.Context, userId int) ([]web.Orders, error)
	GetOrdersItems(ctx context.Context, userId int) ([]web.OrdersItems, error)
	CreateOrders(ctx context.Context, ord web.Orders) (int, error)
	CreateOrdersItems(ctx context.Context, ord web.OrdersItems) error
	UpdateOrders(ctx context.Context, ord web.Orders) error
}

type ServiceImpl struct {
	repo Repository
	db   *sql.DB
}

func NewService(repo Repository, db *sql.DB) Service {
	return &ServiceImpl{
		repo: repo,
		db:   db,
	}
}

func (s *ServiceImpl) GetOrders(ctx context.Context, userId int) ([]web.Orders, error) {

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	ord, err := s.repo.GetOrders(ctx, tx, userId)
	if err != nil {
		return nil, err
	}

	return utils.ConvertOrdersToSlice(ord), nil
}

func (s *ServiceImpl) GetOrdersItems(ctx context.Context, userId int) ([]web.OrdersItems, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	ord, err := s.repo.GetOrdersItems(ctx, tx, userId)
	if err != nil {
		return nil, err
	}
	return utils.ConvertItemsToSlice(ord), nil
}

func (s *ServiceImpl) CreateOrders(ctx context.Context, ord web.Orders) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer utils.CommitOrRollback(tx)

	prd := &domain.Orders{
		Id:      ord.Id,
		Userid:  ord.Userid,
		Status:  "pending",
		Address: ord.Address,
	}

	id, err := s.repo.CreateOrders(ctx, tx, *prd)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *ServiceImpl) CreateOrdersItems(ctx context.Context, ord web.OrdersItems) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	err = s.repo.CreateOrdersItems(ctx, tx, domain.OrdersItems{
		Id:       ord.Id,
		Userid:   ord.Userid,
		Produkid: ord.Produkid,
		Orderid:  ord.Orderid,
		Quantity: ord.Quantity,
		Price:    ord.Price,
	})
	if err != nil {
		return err
	}
	return nil

}

func (s *ServiceImpl) UpdateOrders(ctx context.Context, ord web.Orders) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	prd := &domain.Orders{
		Id:      ord.Id,
		Userid:  ord.Userid,
		Status:  "sukses",
		Address: ord.Address,
	}

	err = s.repo.UpdateOrders(ctx, tx, *prd)
	return err
}
