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
	server       *Server
	conn         *websocket.Conn
	ctx          context.Context
	stop         context.CancelFunc
	onDisconnect ClientFunc
	user         database.User
	remoteAddr   string
}

func Open(server *Server, user database.User, serverCtx context.Context, w http.ResponseWriter, r *http.Request, onDisconnect ClientFunc) (*Client, error) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading HTTP request to websocket connection, %s\n", err.Error())
		return nil, err
	}
	clientCtx, cancel := context.WithCancel(serverCtx)
	client := Client{
		server:       server,
		conn:         c,
		ctx:          clientCtx,
		stop:         cancel,
		onDisconnect: onDisconnect,
		user:         user,
		remoteAddr:   r.RemoteAddr,
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
	for {
		var message database.Message
		err = wsjson.Read(client.ctx, client.conn, &message)
		if err != nil {
			break
		}

		if _, err := client.server.db.InsertMessage(message); err != nil {
			log.Printf("Error sending message from user %s, client %s, message %+v", client.user.Username, client.remoteAddr, message)
		} else {
			// Only send to the recipient if the insert succeeded
			err := client.server.RunWithClient(message.Recipient.UUID, func(client *Client) error {
				return client.SendMessage(message)
			})
			if err != nil {
				log.Printf("Error sending realtime message from user %s, message %+v", client.user.Username, message)
			}
		}
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
	_ = client.onDisconnect(client)
}
