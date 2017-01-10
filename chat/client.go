package chat

import (
	"fmt"
	"io"
	"log"

	"../websocket"
)

const channelBufSize = 100

var maxId int = 0

// Chat client.
type Client struct {
	id     int
	userid string
	ws     *websocket.Conn
	server *Server
	ch     chan Message
	doneCh chan bool
}

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{maxId, "", ws, server, ch, doneCh}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg)
			_, err := c.ws.Write(msg)
			if err != nil {
				c.server.Del(c)
				c.doneCh <- true // for listenRead method
				log.Println("Error:", err.Error())
				return
			}

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			/*var msg Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
			} else {
				c.server.SendAll(&msg)
			}
			*/
			var msg Message = make(Message, 128)
			_, err := c.ws.Read(msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				log.Println("Error:", err.Error())
			} else {
				c.server.SendAll(msg)
				input := ParseMessage(msg)
				log.Printf("Receive: %s\n", msg[:])
				log.Printf("%+v", input)
				if input == nil {
					action, loginin, user := WhetherLogin(msg)
					log.Println(action, loginin, user)
				}
			}

		}
	}
}
