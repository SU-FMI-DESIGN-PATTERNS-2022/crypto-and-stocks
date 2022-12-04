package stream

import "time"

type AuthRequest struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

type SubscriptionRequest struct {
	Action string   `json:"action"`
	Quotes []string `json:"quotes"`
}

type Response struct {
	Type     string    `json:"T"`
	Message  string    `json:"msg"`
	Code     int       `json:"code"`
	Symbol   string    `json:"S"`
	BidPrice float64   `json:"bp"`
	BidSize  float64   `json:"bs"`
	AskPrice float64   `json:"ap"`
	AskSize  float64   `json:"as"`
	Date     time.Time `json:"t"`
}

type CryptoResponse struct {
	Response
	Exchange string `json:"x"`
}

type StockResponse struct {
	Response
	AskExchange string   `json:"ax"`
	BidExchange string   `json:"bx"`
	TradeSize   float64  `json:"s"`
	Conditions  []string `json:"c"`
	Tape        string   `json:"z"`
}
