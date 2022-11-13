package prices

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
	Exchange string    `json:"x"`
	BidPrice float64   `json:"bp"`
	AskPrice float64   `json:"ap"`
	Date     time.Time `json:"t"`
}
