package produk

import (
	"context"
	"database/sql"
	"ddd2/model/domain"
	"ddd2/model/web"
	"ddd2/utils"
	"strings"
)

type ServiceImpl struct {
	Repo Repository
	Db   *sql.DB
}

func NewService(repo Repository, db *sql.DB) Service {
	return &ServiceImpl{
		Repo: repo,
		Db:   db,
	}
}

type Service interface {
	GetAllProduk(ctx context.Context) ([]web.Produk, error)
	GetLikeProduk(ctx context.Context, name string) ([]web.Produk, error)
	GetById(ctx context.Context, id []int) ([]web.Produk, error)
	CreateProduk(ctx context.Context, p *web.ProdukCreatePayload) error
	UpdateProduk(ctx context.Context, p *web.Produk) error
	GetProduk(ctx context.Context, name string) ([]web.Produk, error)
}

func (s *ServiceImpl) GetAllProduk(ctx context.Context) ([]web.Produk, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	p, err := s.Repo.GetAllProduk(ctx, tx)
	if err != nil {
		return nil, err
	}

	return utils.ConvertProdukSlice(p), nil
}

func (s *ServiceImpl) GetLikeProduk(ctx context.Context, name string) ([]web.Produk, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	p, err := s.Repo.GetLikeProduk(ctx, tx, name)
	if err != nil {
		return nil, err
	}

	return utils.ConvertProdukSlice(p), nil
}

func (s *ServiceImpl) GetById(ctx context.Context, id []int) ([]web.Produk, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	p, err := s.Repo.GetById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return utils.ConvertProdukSlice(p), nil
}

func (s *ServiceImpl) CreateProduk(ctx context.Context, p *web.ProdukCreatePayload) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	P := &domain.Produk{
		ProdukName: p.ProdukName,
		Deskripsi:  p.Deskripsi,
		Category:   p.Category,
		UserId:     p.UserId,
		Harga:      p.Harga,
		Quantity:   p.Quantity,
	}
	err = s.Repo.CreateProduk(ctx, tx, *P, p.UserId)
	return err
}

func (s *ServiceImpl) UpdateProduk(ctx context.Context, p *web.Produk) error {
	tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	P := &domain.Produk{
		Id:         p.Id,
		ProdukName: p.ProdukName,
		Deskripsi:  p.Deskripsi,
		Category:   p.Category,
		UserId:     p.UserId,
		Harga:      p.Harga,
		Quantity:   p.Quantity,
	}
	err = s.Repo.UpdateProduk(ctx, tx, *P)
	if err != nil {
		return err
	}
	return nil

}

func (s *ServiceImpl) GetProduk(ctx context.Context, name string) ([]web.Produk, error) {
	tx, err := s.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	strName := strings.ReplaceAll(name, "-", " ")
	p, err := s.Repo.GetProduk(ctx, tx, strName)
	if err != nil {
		return nil, err
	}

	return utils.ConvertProdukSlice(p), nil
}
