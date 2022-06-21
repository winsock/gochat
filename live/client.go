package live

import (
	"context"
	"errors"
	"github.com/winsock/gochat/database"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Client struct {
	conn       *websocket.Conn
	ctx        context.Context
	stop       context.CancelFunc
	user       database.User
	remoteAddr string
}

func Open(user database.User, serverCtx context.Context, w http.ResponseWriter, r *http.Request) (*Client, error) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading HTTP request to websocket connection, %s\n", err.Error())
		return nil, err
	}
	clientCtx, cancel := context.WithCancel(serverCtx)
	client := Client{
		conn:       c,
		ctx:        clientCtx,
		stop:       cancel,
		user:       user,
		remoteAddr: r.RemoteAddr,
	}
	go client.serve()
	return &client, nil
}

func (client *Client) SendMessage(message database.Message) error {
	log.Printf("Sending message: %+v\n", message)
	return wsjson.Write(client.ctx, client.conn, message)
}

func (client *Client) Close() {
	client.stop()
}

func (client *Client) serve() {
	var err error
	for err == nil {
		var message database.Message
		err = wsjson.Read(client.ctx, client.conn, &message)
	}
	log.Printf("Closing connection to client %+v\n", client.user)
	if errors.Is(err, context.Canceled) {
		if err := client.conn.Close(websocket.StatusNormalClosure, "server closing connection"); err != nil {
			log.Printf("Error closing connection, %s\n", err.Error())
		}
	} else {
		if err := client.conn.Close(websocket.StatusInternalError, "server closing connection"); err != nil {
			log.Printf("Error closing connection, %s\n", err.Error())
		}
	}
}
