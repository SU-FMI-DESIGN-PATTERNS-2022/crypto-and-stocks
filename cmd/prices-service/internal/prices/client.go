package prices

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type ClientSocetConfig struct {
	URL    string
	Quotes []string
	Key    string
	Secret string
}

type ClientSocket struct {
	conn *websocket.Conn
}

func NewClientSocket(clientSocetConfig ClientSocetConfig) (*ClientSocket, error) {
	conn, err := dial(clientSocetConfig.URL)
	if err != nil {
		return nil, err
	}

	if err := authenticate(conn, clientSocetConfig.Key, clientSocetConfig.Secret); err != nil {
		return nil, err
	}

	if err := subscribe(conn, clientSocetConfig.Quotes); err != nil {
		return nil, err
	}

	return &ClientSocket{conn: conn}, nil
}

func (c *ClientSocket) Read() {
	resp, _ := readResponse(c.conn, "q", "")
	fmt.Println(resp)
}

func dial(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	if _, err := readResponse(conn, "success", "connected"); err != nil {
		return nil, err
	}

	return conn, nil
}

func authenticate(conn *websocket.Conn, key, secret string) error {
	authRequest := AuthRequest{
		Action: "auth",
		Key:    key,
		Secret: secret,
	}

	authRequestJSON, err := json.Marshal(authRequest)
	if err != nil {
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, authRequestJSON); err != nil {
		return err
	}

	if _, err := readResponse(conn, "success", "authenticated"); err != nil {
		return err
	}

	return nil
}

func subscribe(conn *websocket.Conn, quotes []string) error {
	subRequest := SubscriptionRequest{
		Action: "subscribe",
		Quotes: quotes,
	}

	subRequestJSON, err := json.Marshal(subRequest)
	if err != nil {
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, subRequestJSON); err != nil {
		return err
	}

	if _, err := readResponse(conn, "subscription", ""); err != nil {
		return err
	}

	return nil
}

func readResponse(conn *websocket.Conn, expectedType, expectedMsg string) (Response, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return Response{}, err
	}

	var response []Response
	if err := json.Unmarshal(message, &response); err != nil {
		return Response{}, err
	}

	if response[0].Type != expectedType || response[0].Message != expectedMsg {
		return Response{}, fmt.Errorf("%d: %s", response[0].Code, response[0].Message)
	}

	return response[0], nil
}
