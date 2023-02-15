package order

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/account-service/internal/repositories/order_repository"
	"github.com/gorilla/websocket"
)

type OrderHandler struct {
	presenter OrderPresenter
}

type Message struct {
	Message string `json:"message"`
}

var (
	getOrdersBySymbolExpr          = regexp.MustCompile(`^\/orders\/symbol\/(?P<Param>\w+)$`)
	getOrdersByUserExpr            = regexp.MustCompile(`^\/orders\/user\/(?P<Param>\w+)$`)
	getOrdersByUserIdAndSymbolExpr = regexp.MustCompile(`^\/orders\/user\/symbol$`)
	createUserExpr                 = regexp.MustCompile(`^\/create\/user$`)
	createBotExpr                  = regexp.MustCompile(`^\/create\/bot$`)
	mergeUserAndBotExpr            = regexp.MustCompile(`^\/merge$`)
	estimateUserAmountExpr         = regexp.MustCompile(`^\/user\/amount\/(?P<Param>\w+)$`)
)

func NewOrderHandler(orderPresenter OrderPresenter) OrderHandler {
	return OrderHandler{
		presenter: orderPresenter,
	}
}

func (handler *OrderHandler) notFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	message, _ := json.Marshal(Message{"Not found!"})
	res.Write(message)
}

func (handler *OrderHandler) badRequest(res http.ResponseWriter, req *http.Request, msg string) {
	res.WriteHeader(http.StatusBadRequest)
	message, _ := json.Marshal(Message{msg})
	res.Write(message)
}

func (handler *OrderHandler) internalServerError(res http.ResponseWriter, req *http.Request, msg string) {
	res.WriteHeader(http.StatusInternalServerError)
	message, _ := json.Marshal(Message{msg})
	res.Write(message)
}

func (handler *OrderHandler) success(res http.ResponseWriter, req *http.Request, body interface{}) {
	res.WriteHeader(http.StatusOK)
	message, _ := json.Marshal(body)
	res.Write(message)
}

func (handler *OrderHandler) GetAllOrders(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	if req.Method != http.MethodGet {
		handler.notFound(res, req)
		return
	}

	orders, err := handler.presenter.GetAllOrders()
	if err != nil {
		handler.internalServerError(res, req, "Could not fetch orders!")
		return
	}
	handler.success(res, req, orders)
}

func (handler *OrderHandler) GetAllOrdersByUserId(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	if !getOrdersByUserExpr.MatchString(req.URL.Path) {
		handler.notFound(res, req)
		return
	}

	idParam := getOrdersByUserExpr.FindStringSubmatch(req.URL.Path)
	id, err := strconv.ParseInt(idParam[1], 10, 64)
	if err != nil {
		handler.badRequest(res, req, "Parameter id should be number")
		return
	}

	orders, err := handler.presenter.GetAllOrdersByUserId(id)
	if err != nil {
		handler.internalServerError(res, req, "Could not fetch orders!")
		return
	}
	handler.success(res, req, orders)
}

func (handler *OrderHandler) GetAllOrdersByUserIdAndSymbol(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	if !getOrdersByUserIdAndSymbolExpr.MatchString(req.URL.Path) {
		handler.notFound(res, req)
		return
	}
	query := req.URL.Query()
	if query.Get("id") == "" {
		handler.badRequest(res, req, "Missing query parameter 'id'")
		return
	}
	if query.Get("symbol") == "" {
		handler.badRequest(res, req, "Missing query parameter 'symbol'")
		return
	}

	id, err := strconv.ParseInt(query.Get("id"), 10, 64)
	if err != nil {
		handler.badRequest(res, req, "Query parameter id should be number")
		return
	}

	orders, err := handler.presenter.GetAllOrdersByUserIdAndSymbol(id, query.Get("symbol"))
	if err != nil {
		handler.internalServerError(res, req, "Could not fetch orders!")
		return
	}
	handler.success(res, req, orders)
}

func (handler *OrderHandler) GetAllOrdersBySymbol(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	if !getOrdersBySymbolExpr.MatchString(req.URL.Path) {
		handler.notFound(res, req)
		return
	}
	symbol := getOrdersBySymbolExpr.FindStringSubmatch(req.URL.Path)[1]

	orders, err := handler.presenter.GetAllOrdersBySymbol(symbol)
	if err != nil {
		handler.internalServerError(res, req, "Could not fetch orders!")
		return
	}
	handler.success(res, req, orders)
}

func (handler *OrderHandler) CreateUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	if !createUserExpr.MatchString(req.URL.Path) {
		handler.notFound(res, req)
		return
	}

	query := req.URL.Query()
	if query.Get("id") == "" {
		handler.badRequest(res, req, "Missing query parameter 'id'")
		return
	}
	if query.Get("name") == "" {
		handler.badRequest(res, req, "Missing query parameter 'name'")
		return
	}

	id, err := strconv.ParseInt(query.Get("id"), 10, 64)
	if err != nil {
		handler.badRequest(res, req, "Query parameter id should be number")
		return
	}

	reqErr := handler.presenter.CreateUser(id, query.Get("name"))
	if reqErr != nil {
		handler.badRequest(res, req, reqErr.Error())
		return
	}
	handler.success(res, req, Message{"User created successfully!"})
}

func (handler *OrderHandler) CreateBot(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	if !createBotExpr.MatchString(req.URL.Path) {
		handler.notFound(res, req)
		return
	}

	query := req.URL.Query()
	if query.Get("id") == "" {
		handler.badRequest(res, req, "Missing query parameter 'id'")
		return
	}
	if query.Get("amount") == "" {
		handler.badRequest(res, req, "Missing query parameter 'amount'")
		return
	}

	id, err := strconv.ParseInt(query.Get("id"), 10, 64)
	if err != nil {
		handler.badRequest(res, req, "Query parameter id should be number")
		return
	}

	amount, err := strconv.ParseFloat(query.Get("amount"), 64)
	if err != nil {
		handler.badRequest(res, req, "Query parameter amount should be number")
		return
	}

	reqErr := handler.presenter.CreateBot(id, amount)
	if reqErr != nil {
		handler.badRequest(res, req, reqErr.Error())
		return
	}
	handler.success(res, req, Message{"Bot successfully created!"})
}

func (handler *OrderHandler) MergeUserAndBot(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	if !mergeUserAndBotExpr.MatchString(req.URL.Path) {
		handler.notFound(res, req)
		return
	}

	query := req.URL.Query()
	if query.Get("id") == "" {
		handler.badRequest(res, req, "Missing query parameter 'id'")
		return
	}

	id, err := strconv.ParseInt(query.Get("id"), 10, 64)
	if err != nil {
		handler.badRequest(res, req, "Query parameter id should be number")
		return
	}

	reqErr := handler.presenter.MergeUserAndBot(id)
	if reqErr != nil {
		handler.badRequest(res, req, reqErr.Error())
		return
	}
	handler.success(res, req, Message{"Successfully merged user and bot"})
}

func (handler *OrderHandler) EstimateUserAmount(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-type", "application/json")
	if !estimateUserAmountExpr.MatchString(req.URL.Path) {
		handler.notFound(res, req)
		return
	}

	idParam := estimateUserAmountExpr.FindStringSubmatch(req.URL.Path)
	id, err := strconv.ParseInt(idParam[1], 10, 64)
	if err != nil {
		handler.badRequest(res, req, "Parameter id should be number")
		return
	}
	amount, err := handler.presenter.EstimateUserAmount(id)
	if err != nil {
		handler.internalServerError(res, req, "internal server error")
		return
	}
	handler.success(res, req, map[string]float64{"amount": amount})
}

func (handler *OrderHandler) StoreOrder(res http.ResponseWriter, req *http.Request) {
	conn, err := handler.presenter.upgrader.Upgrade(res, req, nil)

	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage() // json format
		if err != nil {
			fmt.Println("Failed to read message:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Hello, something is wrong."))
			break
		}
		// deserialize to struct Order
		var order order_repository.Order
		if err := json.Unmarshal(message, &order); err != nil {
			fmt.Println(string(message))
			fmt.Println("Failed to unmarshal message:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("The message is not in right json object structure."))
			break
		}
		if err := handler.presenter.StoreOrder(order); err != nil {
			fmt.Println("Failed to store order:", err)
			conn.WriteMessage(websocket.TextMessage, []byte("We have a problem with storing your order."))
			break
		}
	}
}
