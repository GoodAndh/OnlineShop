package cart

import (
	"context"
	"ddd2/exception"
	"ddd2/layer/auth"
	"ddd2/layer/produk"
	"ddd2/model/web"
	"ddd2/utils"
	"ddd2/views"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	service  Service
	produk   produk.Service
	validate *validator.Validate
}

func NewRoute(service Service, produk produk.Service, validate *validator.Validate) *Route {
	return &Route{
		service:  service,
		produk:   produk,
		validate: validate,
	}
}

func (h *Route) RegisterRoute(router *httprouter.Router) {
	router.POST("/p/g=:produkname", h.handleAddToCart)

	router.GET("/p/cart", h.handleCheckout)
	router.POST("/p/cart", h.handleCheckout)

}

func (h *Route) handleAddToCart(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		pnme := params.ByName("produkname")
		errorList := make(map[any]any)
		ses, err := auth.Session.Get(r, "lg-ses")
		if err != nil {
			errorList["error"] = exception.ErrNoSesFound
			exception.WriteInternalServerError(w, errorList, "InfoProduk.html")
			return
		}

		if auten, ok := ses.Values["auten"].(bool); !auten || !ok {
			errorList["error"] = exception.ErrNoSesFound
			exception.WriteInternalServerError(w, errorList, "InfoProduk.html")
			return
		}

		//get user id
		userId, ok := ses.Values["data"].(int)
		if !ok {
			errorList["error"] = exception.ErrIdNotFound
			exception.WriteInternalServerError(w, errorList, "InfoProduk.html")
			return
		}

		ps, err := h.produk.GetProduk(r.Context(), pnme)
		if err != nil {
			errorList["error"] = err
			exception.WriteInternalServerError(w, errorList, "InfoProduk.html")
			return
		}

		//get produk id
		var pId int
		for _, v := range ps {
			pId = v.Id
		}

		qty, _ := strconv.Atoi(r.Form.Get("jumlah"))

		pItems := &web.CartItem{
			Produkid: pId,
			Quantity: qty,
		}
		payload := &web.CartCheckoutPayload{}
		payload.Items = append(payload.Items, *pItems)

		errorList, err = utils.ValidationField(h.validate, *payload)
		if err != nil {
			errorList["error"] = err
			exception.WriteInternalServerError(w, errorList, "InfoProduk.html")
			return
		}

		err = h.service.CreateOrders(context.Background(), ps, payload.Items, userId)
		if err != nil {
			errorList["error"] = err
			exception.WriteInternalServerError(w, errorList, "InfoProduk.html")
			return
		}

	}
}

func (h *Route) handleCheckout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch r.Method {
	case http.MethodGet:
		errorList := make(map[any]any)

		ses, err := auth.Session.Get(r, "lg-ses")
		if err != nil {
			errorList["error"] = exception.ErrNoSesFound
			exception.WriteInternalServerError(w, errorList, "cart.html")
			return
		}
		if auten, ok := ses.Values["auten"].(bool); !auten || !ok {
			errorList["error"] = exception.ErrNoSesFound

			exception.WriteInternalServerError(w, errorList, "cart.html")

			return
		}

		//get user id
		userId, ok := ses.Values["data"].(int)
		if !ok {
			errorList["error"] = exception.ErrNoSesFound
			exception.WriteInternalServerError(w, errorList, "cart.html")
			return
		}

		//get order id
		orderId, err := h.service.GetOrderId(r.Context(), userId)
		if err != nil {
			errorList["error"] = err
			exception.WriteInternalServerError(w, errorList, "cart.html")
			return
		}

		cart, err := h.service.GetOrdersList(r.Context(), userId, orderId)
		if err != nil {
			errorList["error"] = err
			exception.WriteInternalServerError(w, errorList, "cart.html")
			return
		}
		errorList = cartList(cart)
		for _, v := range cart {
			if v.SemuaHarga >= 0 {
				errorList["semuaharga"] = v.SemuaHarga
			}
		}
		views.TemplateExecuted(w, errorList, "cart.html")
	case http.MethodPost:

		errorList := make(map[any]any)

		r.ParseForm()
		ses, err := auth.Session.Get(r, "lg-ses")
		if err != nil {
			errorList["error"] = exception.ErrNoSesFound
			exception.WriteInternalServerError(w, errorList, "cart.html")
			return
		}

		//get user id
		userId, ok := ses.Values["data"].(int)
		if !ok {
			errorList["error"] = exception.ErrNoSesFound
			exception.WriteInternalServerError(w, errorList, "cart.html")
			return
		}

		//get order id
		orderId, err := h.service.GetOrderId(context.Background(), userId)
		if err != nil {
			errorList["error"] = err
			exception.WriteInternalServerError(w, errorList, "cart.html")
			return
		}

		//get info about order
		cart, err := h.service.GetOrdersList(context.Background(), userId, orderId)
		if err != nil {
			errorList["error"] = err
			exception.WriteInternalServerError(w, errorList, "cart.html")
			return
		}

		//get produk id 
		pId := make([]int, len(cart))
		for _, v := range cart {
			
			
			pId = append(pId, v.Produkid)
		}

		//get the product stock left
		ps, err := h.produk.GetById(r.Context(), pId)
		if err != nil {
			errorList["error"] = err
			exception.WriteInternalServerError(w, errorList, "cart.html")
			return
		}

		//store ps into map for easier acces
		pMap := make(map[int]web.Produk)
		for _, items := range ps {
			pMap[items.Id] = items
		}

		log.Println("sebelum berubah = ", pMap)
		//reduce the quantity of product in database
		for _, items := range cart {
			if items.Produkid >=1{
				product := pMap[items.Produkid]
				product.Quantity -= items.Quantity
				err = h.produk.UpdateProduk(r.Context(), &product)
				if err != nil {
					errorList["error"] = err
					exception.WriteInternalServerError(w, errorList, "cart.html")
					return
				}
			}
		}

		//update orders status
		for _,id  := range orderId {
			err=h.service.UpdateOrders(r.Context(),web.Orders{
				Id:      id,
				Status:  "sukses",
				Address: "some address",
			})
			if err != nil {
				errorList["error"] = err
				exception.WriteInternalServerError(w, errorList, "cart.html")
				return
			}

		}



		log.Println("sesudah berubah = ", pMap)

	}

}

func cartList(cart []CartService) map[any]any {
	l1 := make(map[any]any)
	l := make(map[any]any)
	for i, v := range cart {
		if v.Orderid >= 1 {
			li := map[any]any{
				"nama":          v.NameProduk,
				"jumlahpesanan": v.Quantity,
				"harga":         v.Price,
				"totalharga":    v.TotalPrice,
			}
			l[strconv.Itoa(i+1)] = li

		}
	}
	l1["data"] = l

	return l1
}
