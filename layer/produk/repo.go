package produk

import (
	"context"
	"database/sql"
	"ddd2/exception"
	"ddd2/model/domain"
	"fmt"
	"strings"
)

type Repositoryimpl struct {
}

type Repository interface {
	GetAllProduk(ctx context.Context, tx *sql.Tx) ([]domain.Produk, error)
	GetLikeProduk(ctx context.Context, tx *sql.Tx, name string) ([]domain.Produk, error)
	GetById(ctx context.Context, tx *sql.Tx, id []int) ([]domain.Produk, error)
	CreateProduk(ctx context.Context, tx *sql.Tx, p domain.Produk, id int) error
	UpdateProduk(ctx context.Context, tx *sql.Tx, pr domain.Produk) error
	GetProduk(ctx context.Context, tx *sql.Tx, name string) ([]domain.Produk, error)
}

func NewRepository() Repository {
	return &Repositoryimpl{}
}

func (r *Repositoryimpl) GetAllProduk(ctx context.Context, tx *sql.Tx) ([]domain.Produk, error) {
	rows, err := tx.QueryContext(ctx, "select * from produk")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	p, err := rowsScanProduk(rows)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *Repositoryimpl) GetLikeProduk(ctx context.Context, tx *sql.Tx, name string) ([]domain.Produk, error) {
	script := "select * from produk where produkname in (select produkname from produk where produkname like ?) or category in (select category from produk where category like ? );"
	rows, err := tx.QueryContext(ctx, script, name+"%", name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	p, err := rowsScanProduk(rows)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *Repositoryimpl) GetById(ctx context.Context, tx *sql.Tx, id []int) ([]domain.Produk, error) {
	placeholders := strings.Repeat("?,", len(id))
	query := fmt.Sprintf("select * from produk where id in (%s%v)", placeholders, 0)

	//convert id to interface
	args := make([]interface{}, len(id))
	for i, v := range id {
		args[i] = v
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	p, err := rowsScanProduk(rows)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repositoryimpl) CreateProduk(ctx context.Context, tx *sql.Tx, p domain.Produk, id int) error {

	err := tx.QueryRowContext(ctx, "select id from user where id = ?", id).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return exception.ErrIdNotFound
		}
		return err
	}

	script := "insert into produk (produkname,deskripsi,category,userid,harga,quantity) values(?,?,?,?,?,?)"
	result, err := tx.ExecContext(ctx, script, p.ProdukName, p.Deskripsi, p.Category, p.UserId, p.Harga, p.Quantity)
	if err != nil {
		return err
	}

	if num, err := result.RowsAffected(); num == 0 || err != nil {
		return fmt.Errorf("cannot find the specific id got %v", p.UserId)
	}

	return nil

}

func (r *Repositoryimpl) UpdateProduk(ctx context.Context, tx *sql.Tx, pr domain.Produk) error {
	result, err := tx.ExecContext(ctx, "update produk set produkname = ? , deskripsi = ? , category = ? , harga = ? , quantity = ?  where id = ? ", pr.ProdukName, pr.Deskripsi, pr.Category, pr.Harga, pr.Quantity, pr.Id)
	if err != nil {
		return err
	}

	if num, err := result.RowsAffected(); num == 0 || err != nil {
		return fmt.Errorf("cannot find the specific id got %v", pr.Id)
	}

	return nil
}
func (r *Repositoryimpl) GetProduk(ctx context.Context, tx *sql.Tx, name string) ([]domain.Produk, error) {
	rows, err := tx.QueryContext(ctx, "select * from produk where produkname = ? ", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	p, err := rowsScanProduk(rows)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func rowsScanProduk(rows *sql.Rows) ([]domain.Produk, error) {
	var te []domain.Produk
	for rows.Next() {
		var t domain.Produk
		err := rows.Scan(&t.Id, &t.ProdukName, &t.Deskripsi, &t.Category, &t.UserId, &t.Harga, &t.Quantity)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, exception.ErrIdNotFound
			}
			return nil, err
		}
		te = append(te, t)
	}
	return te, nil
}
