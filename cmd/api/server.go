package api

import (
	"database/sql"
	"ddd2/exception"
	"ddd2/layer/cart"
	"ddd2/layer/orders"
	"ddd2/layer/produk"
	"ddd2/layer/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type ApiServer struct {
	Addr     string
	Db       *sql.DB
	Validate *validator.Validate
}

func NewApiServer(addr string, db *sql.DB, validate *validator.Validate) *ApiServer {
	return &ApiServer{
		Addr:     addr,
		Db:       db,
		Validate: validate,
	}
}

func (s *ApiServer) Run() error {

	router := httprouter.New()
	router.NotFound = exception.NewNotFound(exception.NotFound)

	userRepo := user.NewRepository()
	userService := user.NewService(userRepo, s.Db)
	userHandler := user.NewRoute(s.Validate, userService)
	userHandler.RegisterRoute(router)

	produkRepo := produk.NewRepository()
	produkService := produk.NewService(produkRepo, s.Db)
	produkHandler := produk.NewRoute(s.Validate, produkService)
	produkHandler.RegisterRoute(router)

	orderRepo:=orders.NewRepository()
	orderService:=orders.NewService(orderRepo,s.Db)

	cartRepo:=cart.NewRepository()
	cartService:=cart.NewService(orderService,s.Db,cartRepo)
	cartHandler:=cart.NewRoute(cartService,produkService,s.Validate)
	cartHandler.RegisterRoute(router)

	return http.ListenAndServe(s.Addr, router)

}
