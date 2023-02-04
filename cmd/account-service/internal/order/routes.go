package order

import "net/http"

func HandleRoutes(mux *http.ServeMux, handler OrderHandler) {
	mux.HandleFunc("/orders/all", handler.GetAllOrders)
	mux.HandleFunc("/orders/user/", handler.GetAllOrdersByUserId)
	mux.HandleFunc("/orders/user/symbol", handler.GetAllOrdersByUserIdAndSymbol)
	mux.HandleFunc("/orders/symbol/", handler.GetAllOrdersBySymbol)
	mux.HandleFunc("/create/user", handler.CreateUser)
	mux.HandleFunc("/create/bot", handler.CreateBot)
	mux.HandleFunc("/merge", handler.MergeUserAndBot)
	mux.HandleFunc("/user/amount/", handler.EstimateUserAmount)
}
