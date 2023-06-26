package stream

import (
	"encoding/json"
	"fmt"

	"github.com/SU-FMI-DESIGN-PATTERNS-2022/crypto-and-stocks/cmd/prices-service/env"
	"github.com/gorilla/websocket"
)

func NewStreamConfig(wsConfig env.WebSocetConfig) StreamConfig {
	return StreamConfig{
		URL:    wsConfig.CryptoURL,
		Quotes: wsConfig.CryptoQuotes,
		Key:    wsConfig.Key,
		Secret: wsConfig.Secret,
	}
}
func NewCryptoConfig(wsConfig env.WebSocetConfig) StreamConfig {
	return StreamConfig{
		URL:    wsConfig.CryptoURL,
		Quotes: wsConfig.CryptoQuotes,
		Key:    wsConfig.Key,
		Secret: wsConfig.Secret,
	}
}

func NewStockConfig(wsConfig env.WebSocetConfig) StreamConfig {
	return StreamConfig{
		URL:    wsConfig.StockURL,
		Quotes: wsConfig.StockQuotes,
		Key:    wsConfig.Key,
		Secret: wsConfig.Secret,
	}
}

type StreamConfig struct {
	URL    string
	Quotes []string
	Key    string
	Secret string
}

type socetConn struct {
	*websocket.Conn
}

type Stream struct {
	conn        *socetConn
	priceStream chan []byte
	closeChan   chan struct{}
}

func newSocketConnection(url string) (*socetConn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, fmt.Errorf("dial websocket fails: %w", err)
	}

	return &socetConn{conn}, nil
}

func NewPriceStream(clientSocetConfig StreamConfig) (*Stream, error) {
	conn, err := newSocketConnection(clientSocetConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("create new socket connection fails: %w", err)
	}

	if _, err := conn.readResponse("success", "connected"); err != nil {
		return nil, fmt.Errorf("read connection response fails: %w", err)
	}

	if err := conn.authenticate(clientSocetConfig.Key, clientSocetConfig.Secret); err != nil {
		return nil, fmt.Errorf("authenticate to socket fails: %w", err)
	}

	if err := conn.subscribe(clientSocetConfig.Quotes); err != nil {
		return nil, fmt.Errorf("subscibe to quotes fails: %w", err)
	}

	return &Stream{
		conn:        conn,
		priceStream: make(chan []byte, 1),
		closeChan:   make(chan struct{}, 1),
	}, nil
}

func (s *Stream) Start(msgHandler func([]byte) error) error {
	errChan := make(chan error, 1)

	go s.listenForResponse(msgHandler, errChan)

	return <-errChan
}

func (s *Stream) listenForResponse(msgHandler func([]byte) error, errChan chan error) {
	defer func() {
		s.conn.Close()
	}()

	for {
		select {
		case <-s.closeChan:
			return
		default:
			resp, err := s.conn.readResponse("q", "")
			if err != nil {
				errChan <- fmt.Errorf("read quote response fails: %w", err)
				return
			}

			if err := msgHandler(resp); err != nil {
				errChan <- fmt.Errorf("handling message fails: %w", err)
				return
			}
		}
	}
}

func (s *Stream) Stop() {
	close(s.closeChan)
}

func (sc *socetConn) authenticate(key, secret string) error {
	authRequest := AuthRequest{
		Action: "auth",
		Key:    key,
		Secret: secret,
	}

	authRequestJSON, err := json.Marshal(authRequest)
	if err != nil {
		return fmt.Errorf("marshal authenticate request fails: %w", err)
	}

	if err := sc.WriteMessage(websocket.TextMessage, authRequestJSON); err != nil {
		return fmt.Errorf("write authenticate request to socket fails: %w", err)
	}

	if _, err := sc.readResponse("success", "authenticated"); err != nil {
		return fmt.Errorf("read authenticate response fails: %w", err)
	}

	return nil
}

func (sc *socetConn) subscribe(quotes []string) error {
	subRequest := SubscriptionRequest{
		Action: "subscribe",
		Quotes: quotes,
	}

	subRequestJSON, err := json.Marshal(subRequest)
	if err != nil {
		return fmt.Errorf("marshal subscribe request fails: %w", err)
	}

	if err := sc.WriteMessage(websocket.TextMessage, subRequestJSON); err != nil {
		return fmt.Errorf("write subscribe request to socket fails: %w", err)
	}

	if _, err := sc.readResponse("subscription", ""); err != nil {
		return fmt.Errorf("read subscribe response fails: %w", err)
	}

	return nil
}

func (sc *socetConn) readResponse(expectedType, expectedMsg string) ([]byte, error) {
	_, message, err := sc.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read response form socket fails: %w", err)
	}

	var response []Response
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response fails: %w", err)
	}

	for _, responseMsg := range response {
		if responseMsg.Type != expectedType || responseMsg.Message != expectedMsg {
			return nil, fmt.Errorf("unexpected response: %d: %s", responseMsg.Code, responseMsg.Message)
		}
	}

	return message, nil
}
