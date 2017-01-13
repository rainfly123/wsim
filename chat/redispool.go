package chat

import "menteslibres.net/gosexy/redis"
import "log"

var host = "127.0.0.1"
var port = uint(6379)
var clients redisPool

type redisPool struct {
	connections chan *redis.Client
	connFn      func() (*redis.Client, error) // function to create new connection.
}

func (this *redisPool) Get() (*redis.Client, bool) {
	var conn *redis.Client
	select {
	case conn = <-this.connections:
	default:
		conn, err := this.connFn()
		if err != nil {
			return nil, false
		}

		return conn, true
	}

	if err := this.testConn(conn); err != nil {
		return this.Get() // if connection is bad, get the next one in line until base case is hit, then create new client
	}

	return conn, true
}

func (this *redisPool) Close(conn *redis.Client) {
	select {
	case this.connections <- conn:
		return
	default:
		conn.Quit()
	}
}

func (this *redisPool) testConn(conn *redis.Client) error {
	if _, err := conn.Ping(); err != nil {
		conn.Quit()
		return err
	}

	return nil
}
func newcon() (*redis.Client, error) {
	var client *redis.Client
	client = redis.New()
	err := client.Connect(host, port)
	client.Command(nil, "SELECT", 1)
	return client, err
}

func InitRedis() {
	clients.connFn = newcon
	clients.connections = make(chan *redis.Client, 10)
}
func PushOfflineMsg(user string, msg string) {
	var client *redis.Client
	var ok bool

	client, ok = clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.RPush(user, msg)
	client.Close()
}
func SendoutOfflineMsg(user string) []string {
	var client *redis.Client
	var ok bool
	var null []string

	client, ok = clients.Get()
	if ok != true {
		log.Panic("redis error")
		return null
	}
	msg, _ := client.LRange(user, 0, -1)
	client.Del(user)
	client.Close()
	return msg
}
