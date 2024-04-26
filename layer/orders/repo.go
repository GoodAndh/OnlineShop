package orders

import (
	"context"
	"database/sql"
	"ddd2/model/domain"
	"fmt"
)

type Repository interface {
	GetOrders(ctx context.Context, tx *sql.Tx, userId int) ([]domain.Orders, error)
	GetOrdersItems(ctx context.Context, tx *sql.Tx, userId int) ([]domain.OrdersItems, error)
	CreateOrders(ctx context.Context, tx *sql.Tx, ord domain.Orders) (int, error)
	CreateOrdersItems(ctx context.Context, tx *sql.Tx, ord domain.OrdersItems) error
	UpdateOrders(ctx context.Context, tx *sql.Tx, ord domain.Orders) error
}

type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}
func (r *RepositoryImpl) GetOrders(ctx context.Context, tx *sql.Tx, userId int) ([]domain.Orders, error) {
	rows, err := tx.QueryContext(ctx, "select * from orders where userid = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ord := make([]domain.Orders, 0)
	for rows.Next() {
		ORD := domain.Orders{}
		err := rows.Scan(&ORD.Id, &ORD.Userid, &ORD.Status, &ORD.Address)
		if err != nil {
			return nil, err
		}
		ord = append(ord, ORD)

	}
	return ord, nil
}

func (r *RepositoryImpl) GetOrdersItems(ctx context.Context, tx *sql.Tx, userId int) ([]domain.OrdersItems, error) {
	rows, err := tx.QueryContext(ctx, "select * from orders_items where userId = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ord := make([]domain.OrdersItems, 0)
	for rows.Next() {
		ORD := domain.OrdersItems{}
		err := rows.Scan(&ORD.Id, &ORD.Userid, &ORD.Produkid, &ORD.Quantity, &ORD.Price)
		if err != nil {
			return nil, err
		}
		ord = append(ord, ORD)
	}
	return ord, nil
}

func (r *RepositoryImpl) CreateOrders(ctx context.Context, tx *sql.Tx, ord domain.Orders) (int, error) {
	result, err := tx.ExecContext(ctx, "insert into orders (userid,status,address) values(?,?,?)", ord.Userid, ord.Status, ord.Address)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return 0, fmt.Errorf("have error while creating orders rows affected = %v", affected)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (r *RepositoryImpl) CreateOrdersItems(ctx context.Context, tx *sql.Tx, ord domain.OrdersItems) error {
	result, err := tx.ExecContext(ctx, "insert into orders_items(userid,produkid,orderid,quantity,price) values (?,?,?,?,?)", ord.Userid, ord.Produkid, ord.Orderid, ord.Quantity, ord.Price)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return fmt.Errorf("have error while creating orders_items rows affected = %v", affected)
	}
	return nil
}

func (r *RepositoryImpl) UpdateOrders(ctx context.Context, tx *sql.Tx, ord domain.Orders) error {
	result,err:=tx.ExecContext(ctx,"update orders set status = ? ,address = ? where id = ",ord.Status,ord.Address,ord.Id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected == 0 {
		return fmt.Errorf("have error while creating orders_items rows affected = %v", affected)
	}
	return nil
}
