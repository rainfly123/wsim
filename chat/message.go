package chat

type Message []byte

func (self Message) String() string {
	return string(self)
}
