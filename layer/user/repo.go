package user

import (
	"context"
	"database/sql"
	"ddd2/model/domain"
	"log"
)

type Repositoryimpl struct {
}

type Repository interface {
	GetByUsername(ctx context.Context, tx *sql.Tx, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*domain.User, error)
	GetById(ctx context.Context, tx *sql.Tx, id int) (*domain.User, error)
	CreateUser(ctx context.Context, tx *sql.Tx, user domain.User) error
}

func NewRepository() Repository {
	return &Repositoryimpl{}
}

func (r *Repositoryimpl) GetByUsername(ctx context.Context, tx *sql.Tx, username string) (*domain.User, error) {

	rows, err := tx.QueryContext(ctx, "select * from user where username = ? ", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	u, err := rowsScanUser(rows)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repositoryimpl) GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*domain.User, error) {

	rows, err := tx.QueryContext(ctx, "select * from user where username = ? ", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	e, err := rowsScanUser(rows)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *Repositoryimpl) GetById(ctx context.Context, tx *sql.Tx, id int) (*domain.User, error) {
	rows,err:=tx.QueryContext(ctx,"select * from user where id = ; ?",id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	u,err:=rowsScanUser(rows)
	if err != nil {
		return nil, err
	}
	return u,nil
}

func (r *Repositoryimpl) CreateUser(ctx context.Context, tx *sql.Tx, user domain.User) error {

	_, err := tx.ExecContext(ctx, "insert into user (username,password,email,name) values(?,?,?,?)", user.Username, user.Password, user.Email, user.Name)
	if err != nil {
		return err
	}
	log.Println("u adalah=", user)

	return nil

}

func rowsScanUser(rows *sql.Rows) (*domain.User, error) {
	u := &domain.User{}
	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Username, &u.Password, &u.Name, &u.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, err
			}
			return nil, err
		}
	}
	return u, nil
}
