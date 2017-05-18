package gateway

const LocalIPRecvAddr = "192.168.1.240:2222"
const LocalIPSendAddr = "192.168.1.240:1111"

type Transmit struct {
	RemoteIPAddr string
	Userid       string
	Message      string
}

//redis
const host = "127.0.0.1"
const port = uint(6379)
