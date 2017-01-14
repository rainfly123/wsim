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
			log.Println("Send:", msg, c.userid)
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
			var msg Message = make(Message, 384)
			i, err := c.ws.Read(msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				log.Println("Error:", err.Error())
			} else {
				//c.server.SendAll(msg)
				msg = msg[:i]
				input := ParseMessage(msg, c.userid)
				log.Printf("Receive: %s\n", msg[:])
				if input == nil { //not emotion_/picture_/video_
					action, login, user := WhetherLogin(msg)
					log.Println(action, login, user)
					if action {
						if login {
							c.userid = user
							c.server.Online(c)
							messages := SendoutOfflineMsg(user)
							for _, v := range messages {
								c.Write([]byte(v)) //offline message
							}
						} else {
							c.server.Del(c)
						}
					}
				} else { //emotion_/picture_/video_

					if input.Ttype == GRP {
						//group cast
						GroupMsgCh <- input
					} else {
						//unicast
						lockUsers.RLock()
						touser, online := c.server.users[input.Touserid]
						lockUsers.RUnlock()
						output := NewOutput(input)
						if online {
							touser.Write([]byte(output.String()))
						} else {
							//offline....
							PushOfflineMsg(input.Touserid, output.String())
						}
					}
				}
			}

		}
	}
}
