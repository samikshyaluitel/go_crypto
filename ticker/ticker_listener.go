package ticker

import (
	"context"
	"encoding/json"
	"go_crypto/models"

	"github.com/gorilla/websocket"
	"github.com/juju/errors"
	jsonrpc2 "github.com/sourcegraph/jsonrpc2"
	jsonrpc2ws "github.com/sourcegraph/jsonrpc2/websocket"
)

const wsAPIURL string = "wss://api.hitbtc.com/api/2/ws"

// responseChannels handles all incoming data from the hitbtc connection.
type responseChannels struct {
	notifications notificationChannels
	ErrorFeed     chan error
}

// notificationChannels contains all the notifications from hitbtc for subscribed feeds.
type notificationChannels struct {
	TickerFeed map[string]chan models.Ticker
}

// Handle handles all incoming connections and fills the channels properly.
func (h *responseChannels) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if req.Params != nil {
		message := *req.Params
		switch req.Method {
		case "ticker":
			var msg models.Ticker
			err := json.Unmarshal(message, &msg)
			if err != nil {
				h.ErrorFeed <- err
			} else {
				h.notifications.TickerFeed[msg.Symbol] <- msg
			}
		}
	}
}

// TickerClient represents a JSON RPC v2 Connection over Websocket,
type TickerClient struct {
	conn    *jsonrpc2.Conn
	updates *responseChannels
}

// NewTickerClient creates a new TickerClient
func NewTickerClient() (*TickerClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsAPIURL, nil)
	if err != nil {
		return nil, err
	}

	handler := responseChannels{
		notifications: notificationChannels{
			TickerFeed: make(map[string]chan models.Ticker),
		},
		ErrorFeed: make(chan error),
	}

	return &TickerClient{
		conn:    jsonrpc2.NewConn(context.Background(), jsonrpc2ws.NewObjectStream(conn), jsonrpc2.AsyncHandler(&handler)),
		updates: &handler,
	}, nil
}

// Close closes the Websocket connected to the hitbtc api.
func (c *TickerClient) Close() {
	c.conn.Close()

	for _, channel := range c.updates.notifications.TickerFeed {
		close(channel)
	}
	close(c.updates.ErrorFeed)

	c.updates.notifications.TickerFeed = make(map[string]chan models.Ticker)
	c.updates.ErrorFeed = make(chan error)
}

// wsSubscriptionResponse is the response for a subscribe/unsubscribe requests.
type wsSubscriptionResponse bool

// WSSubscriptionRequest is request type on websocket subscription.
type WSSubscriptionRequest struct {
	Symbol string `json:"symbol,required"`
}

// SubscribeTicker subscribes to the specified market ticker notifications.
func (c *TickerClient) SubscribeTicker(symbol string) (<-chan models.Ticker, error) {
	err := c.subscriptionOp("subscribeTicker", symbol)
	if err != nil {
		return nil, errors.Annotate(err, "Hitbtc SubscribeTicker")
	}

	if c.updates.notifications.TickerFeed[symbol] == nil {
		c.updates.notifications.TickerFeed[symbol] = make(chan models.Ticker)
	}

	return c.updates.notifications.TickerFeed[symbol], nil
}

// UnsubscribeTicker subscribes to the specified market ticker notifications.
//
// This closes also the connected channel of updates.
func (c *TickerClient) UnsubscribeTicker(symbol string) error {
	err := c.subscriptionOp("unsubscribeTicker", symbol)
	if err != nil {
		return errors.Annotate(err, "Hitbtc UnsubscribeTicker")
	}

	close(c.updates.notifications.TickerFeed[symbol])
	delete(c.updates.notifications.TickerFeed, symbol)

	return nil
}

const (
	// Interval30Minutes is 30 minutes interval for candle data.
	Interval30Minutes string = "M30"
	// Interval1Hour is 1 hour interval for candle data.
	Interval1Hour string = "H1"
)

func (c *TickerClient) subscriptionOp(op string, symbol string) error {
	if c.conn == nil {
		return errors.New("Connection is unitialized")
	}

	var request = WSSubscriptionRequest{Symbol: symbol}
	var success wsSubscriptionResponse

	err := c.conn.Call(context.Background(), op, request, &success)
	if err != nil {
		return err
	}

	if !success {
		return errors.New("Subscribe not successful")
	}

	return nil
}
