package gateway

const LocalIPRecvAddr = "10.116.129.117:2222"
const LocalIPSendAddr = "10.116.129.117:1111"

type Transmit struct {
	RemoteIPAddr string
	Userid       string
	Message      string
}

//redis
const host = "10.169.205.202"
const port = uint(6379)
