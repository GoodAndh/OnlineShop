package produk

import (
	"ddd2/exception"
	"ddd2/layer/auth"
	"ddd2/model/web"
	"ddd2/utils"
	"ddd2/views"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Validate *validator.Validate
	Service  Service
}

func NewRoute(validate *validator.Validate, service Service) *Route {
	return &Route{
		Validate: validate,
		Service:  service,
	}
}

func (h *Route) RegisterRoute(router *httprouter.Router) {
	router.GET("/", h.handleProduct)
	router.POST("/", h.handleProduct)
	router.GET("/cari", h.handleLike)

	router.GET("/p/c=:create", h.handleCreateProduk)
	router.POST("/p/c=:create", h.handleCreateProduk)

	router.GET("/p/g=:produkname", h.handleProductInfo)
}

func (h *Route) handleProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch r.Method {
	case http.MethodGet:
		Data:=map[any]any{}

		p, err := h.Service.GetAllProduk(r.Context())
		if err != nil {
			Data["err"]=err
			exception.WriteInternalServerError(w, Data, "index.html")
			return
		}

		Data = forRangeProduk(p)
		err = views.TemplateExecuted(w, Data, "index.html")
		if err != nil {
			Data["err"]=err
			exception.WriteInternalServerError(w, Data, "index.html")			
			return
		}
	case http.MethodPost:
		r.ParseForm()

		nm := r.Form.Get("namaproduk")

		time.Sleep(1 * time.Second)
		http.Redirect(w, r, "/cari?cari="+nm, http.StatusSeeOther)

	}
}

func (h *Route) handleLike(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch r.Method {
	case http.MethodGet:
		Data:=make(map[any]any)

		r.ParseForm()

		cari := r.URL.Query().Get("cari")

		p, err := h.Service.GetLikeProduk(r.Context(), cari)
		if err != nil {
			Data["err"]=err
			exception.WriteInternalServerError(w, Data, "index.html")
		}
		Data = forRangeProduk(p)
		err = views.TemplateExecuted(w, Data, "index.html")
		if err != nil {
			Data["err"]=err
			exception.WriteInternalServerError(w, Data, "index.html")
		}

	}
}

func (h *Route) handleCreateProduk(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch r.Method {
	case http.MethodGet:
		errorList:=make(map[any]any)

		rr := params.ByName("create")
		if rr != "create" {
			exception.NotFound(w, r, nil)
			return
		}

		ses, err := auth.Session.Get(r, "lg-ses")
		if err != nil {
			errorList["err"]=exception.ErrNoSesFound
			exception.WriteInternalServerError(w,errorList, "createp.html")
			return
		}

		if auten, ok := ses.Values["auten"].(bool); !auten || !ok {
			errorList["err"]=exception.ErrNoSesFound
			exception.WriteInternalServerError(w,errorList, "createp.html")
			return
		}

		if _, ok := ses.Values["data"].(int);  !ok {
			errorList["err"]=exception.ErrNoSesFound
			exception.WriteInternalServerError(w,errorList, "createp.html")
			return
		}

		err = views.TemplateExecuted(w, nil, "createp.html")
		if err != nil {
			errorList["err"]=exception.ErrNoSesFound
			exception.WriteInternalServerError(w,errorList, "createp.html")
			return
		}
	case http.MethodPost:
		errorList:=make(map[any]any)

		rr := params.ByName("create")
		if rr != "create" {
			exception.NotFound(w, r, nil)
			return
		}

		r.ParseForm()

		ses, err := auth.Session.Get(r, "lg-ses")
		if err != nil {
			errorList["err"]=exception.ErrNoSesFound
			exception.WriteInternalServerError(w,errorList, "createp.html")
			return
		}

		if auten, ok := ses.Values["auten"].(bool); !auten || !ok {
			errorList["err"]=exception.ErrNoSesFound
			exception.WriteInternalServerError(w,errorList, "createp.html")
			return
		}

		value, ok := ses.Values["data"].([]interface{})
		if value == nil || !ok {
			errorList["err"]=exception.ErrNoSesFound
			exception.WriteInternalServerError(w,errorList, "createp.html")
			return
		}
		var pUserid int
		for _, v := range value {
			p := v.([]int)
			for _, val := range p {
				pUserid = val
			}
		}
		pHarga, _ := strconv.Atoi(r.Form.Get("harga"))
		pQuantity, _ := strconv.Atoi(r.Form.Get("quantity"))

		p := &web.ProdukCreatePayload{
			ProdukName: r.Form.Get("nama"),
			Deskripsi:  r.Form.Get("deskripsi"),
			Category:   r.FormValue("category"),
			UserId:     pUserid,
			Harga:      pHarga,
			Quantity:   pQuantity,
		}

		errorList, err = utils.ValidationField(h.Validate, p)
		//bug validate bagian number dan gt not fixed
		if err != nil {
			views.TemplateExecuted(w, errorList, "createp.html")
			return
		}

		if p.Quantity == 0 || p.Harga == 0 || p.UserId == 0 {
			errorList["error"] = "tidak boleh memasukkan angka 0"
			views.TemplateExecuted(w, errorList, "createp.html")
			return
		}

		err = h.Service.CreateProduk(r.Context(), p)
		if err != nil {
			errorList["error"] = err.Error()
			views.TemplateExecuted(w, errorList, "createp.html")
			return
		}
		views.TemplateExecuted(w, errorList, "createp.html")

	}
}

func (h *Route) handleProductInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	switch r.Method {
	case http.MethodGet:

		pn := params.ByName("produkname")

		errorList := make(map[any]any)
		product, err := h.Service.GetProduk(r.Context(), pn)
		if err != nil {
			errorList["error"] = err.Error()
			views.TemplateExecuted(w, errorList, "InfoProduk.html")
			return
		}

		if product == nil {
			errorList["error"] = exception.ErrNotFoundProduk
			views.TemplateExecuted(w, errorList, "InfoProduk.html")
			return
		}

		if len(product) == 0 {
			errorList["error"] = exception.ErrNotFoundProduk
			err := views.TemplateExecuted(w, errorList, "InfoProduk.html")
			if err != nil {
				panic(err)
			}

			return
		}

		errorList = forRangeProduk(product)
		err = views.TemplateExecuted(w, errorList, "InfoProduk.html")
		if err != nil {
			panic(err)
		}

	}
}

func forRangeProduk(p []web.Produk) map[any]any {
	l1 := make(map[any]any)
	l := make(map[any]any)
	for i, v := range p {
		parcit := strings.ReplaceAll(v.ProdukName, " ", "-")
		li := map[any]any{
			"nama":      v.ProdukName,
			"deskripsi": v.Deskripsi,
			"quantity":  v.Quantity,
			"harga":     v.Harga,
			"category":  v.Category,
			"button":    true,
			"namap":     parcit,
		}
		l[strconv.Itoa(i+1)] = li
	}
	l1["data"] = l
	return l1
}
