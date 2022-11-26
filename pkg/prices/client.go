package prices

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

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
		return nil, err
	}

	return &socetConn{conn}, nil
}

func NewStream(clientSocetConfig StreamConfig) (*Stream, error) {
	conn, err := newSocketConnection(clientSocetConfig.URL)
	if err != nil {
		return nil, err
	}

	if _, err := conn.readResponse("success", "connected"); err != nil {
		return nil, err
	}

	if err := conn.authenticate(clientSocetConfig.Key, clientSocetConfig.Secret); err != nil {
		return nil, err
	}

	if err := conn.subscribe(clientSocetConfig.Quotes); err != nil {
		return nil, err
	}

	return &Stream{
		conn:        conn,
		priceStream: make(chan []byte, 1),
		closeChan:   make(chan struct{}, 1),
	}, nil
}

func (s *Stream) Start(msgHandler func([]byte)) error {
	errChan := make(chan error, 1)

	go s.listenForResponse(errChan)
	go s.handleResponse(msgHandler)

	return <-errChan
}

func (s *Stream) listenForResponse(errChan chan error) {
	defer func() {
		close(s.closeChan)
		s.conn.Close()
	}()

	for {
		select {
		case <-s.closeChan:
			return
		default:
			resp, err := s.conn.readResponse("q", "")
			if err != nil {
				errChan <- err
				return
			}

			s.priceStream <- resp
		}
	}
}

func (s *Stream) handleResponse(msgHandler func([]byte)) {
	defer func() {
		close(s.priceStream)
	}()

	for {
		select {
		case <-s.closeChan:
			return
		case response := <-s.priceStream:
			msgHandler(response)
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
		return err
	}

	if err := sc.WriteMessage(websocket.TextMessage, authRequestJSON); err != nil {
		return err
	}

	if _, err := sc.readResponse("success", "authenticated"); err != nil {
		return err
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
		return err
	}

	if err := sc.WriteMessage(websocket.TextMessage, subRequestJSON); err != nil {
		return err
	}

	if _, err := sc.readResponse("subscription", ""); err != nil {
		return err
	}

	return nil
}

func (sc *socetConn) readResponse(expectedType, expectedMsg string) ([]byte, error) {
	_, message, err := sc.ReadMessage()
	if err != nil {
		return nil, err
	}

	var response []Response
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, err
	}

	if response[0].Type != expectedType || response[0].Message != expectedMsg {
		return nil, fmt.Errorf("%d: %s", response[0].Code, response[0].Message)
	}

	return message, nil
}
