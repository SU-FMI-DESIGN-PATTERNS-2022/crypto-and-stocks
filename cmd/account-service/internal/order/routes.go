package order

import (
	"net/http"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/pkg/middleware"
)

func HandleRoutes(mux *http.ServeMux, handler OrderHandler) {
	mux.HandleFunc("/orders/all", middleware.SetContentTypeJSON(handler.GetAllOrders))
	mux.HandleFunc("/orders/user/", middleware.SetContentTypeJSON(handler.GetAllOrdersByUserId))
	mux.HandleFunc("/orders/user/symbol", middleware.SetContentTypeJSON(handler.GetAllOrdersByUserIdAndSymbol))
	mux.HandleFunc("/orders/symbol/", middleware.SetContentTypeJSON(handler.GetAllOrdersBySymbol))
	mux.HandleFunc("/create/user", middleware.SetContentTypeJSON(handler.CreateUser))
	mux.HandleFunc("/create/bot", middleware.Authenticate(middleware.SetContentTypeJSON(handler.CreateBot)))
	mux.HandleFunc("/merge", middleware.Authenticate(middleware.SetContentTypeJSON(handler.MergeUserAndBot)))
	mux.HandleFunc("/user/amount/", middleware.Authenticate(middleware.SetContentTypeJSON(handler.EstimateUserAmount)))
	mux.HandleFunc("/order", middleware.Authenticate(handler.StoreOrder))
}
