package chat

import (
	"log"
	"net/http"
	"sync"

	"../websocket"
)

var lockUsers sync.RWMutex

// Chat server.
type Server struct {
	pattern string
	//	messages  []Message
	//clients   map[int]*Client
	users    map[string]*Client
	addCh    chan *Client
	delCh    chan *Client
	onlineCh chan *Client
	//doneCh   chan bool
	errCh chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	//messages := []Message{}
	//clients := make(map[int]*Client)
	Users := make(map[string]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	onlineCh := make(chan *Client)
	//sendAllCh := make(chan Message)
	//doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		Users,
		addCh,
		delCh,
		onlineCh,
		//doneCh,
		errCh,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Online(c *Client) {
	s.onlineCh <- c
}
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

/*
func (s *Server) SendAll(msg Message) {
	s.sendAllCh <- msg
}
*/
/*
func (s *Server) Done() {
	s.doneCh <- true
}
*/

func (s *Server) Err(err error) {
	s.errCh <- err
}

/*
func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}
*/
/*
func (s *Server) sendAll(msg Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}
*/
// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen() {

	log.Println("Listening server...")

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Server{Handler: onConnected})
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			//log.Println("Added new client")
			//s.clients[c.id] = c
			log.Println("Added new client Now", maxId, "clients.")
			_ = c
		//	s.sendPastMessages(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client", c.userid)
			//delete(s.clients, c.id)
			lockUsers.Lock()
			delete(s.users, c.userid)
			lockUsers.Unlock()

			// broadcast message for all clients
			/*
				case msg := <-s.sendAllCh:
					log.Println("Send all:", msg)
					s.messages = append(s.messages, msg)
					s.sendAll(msg)
			*/
		case c := <-s.onlineCh:
			lockUsers.Lock()
			s.users[c.userid] = c
			lockUsers.Unlock()

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

			//case <-s.doneCh:
			//	return
		}
	}
}

func (s *Server) GetClient(user string) (*Client, bool) {
	lockUsers.RLock()
	client, online := server.users[user]
	lockUsers.RUnlock()
	return client, online
}
