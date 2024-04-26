package exception

import (
	"errors"
)

// err handling user
var (
	ErrIncorrectUser error = errors.New("incorrect username or password")
	ErrUsernameTaken error = errors.New("username telah digunakan")
	ErrEmailTaken    error = errors.New("email telah digunakan")
	ErrIdNotFound error=errors.New("id not found")
)

// err handling produk
var (
	ErrNotFoundProduk error =errors.New("produk yang anda minta belum tersedia")
	ErrQuantityZero error =errors.New("quantity produk equal to zero")
	ErrQuantityMoreThanStock error=errors.New("stock kurang dari permintaan")
)

var (
	ErrCartEmpty error = errors.New("cart is empy")
)

// err handling session
var (
	ErrNoSesFound error = errors.New("no session found ,login needed")
)

// err handling router
var (
	ErrNotFoundRouter error = errors.New("path yang dituju tidak tersedia")
)


