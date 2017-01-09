package chat

import "time"
import "bytes"
import "strings"
import "strconv"

type Message []byte

/*
const emotion = "emotion"
const picture = "picture"
const video = "video"
const unicast = "unicast"
const group = "group"
*/
const (
	EMOTION = 1
	PICTURE = 2
	VIDEO   = 3
)
const (
	UNI = 1
	GRP = 2
)

type OutPut struct {
	Mtype      uint8
	Fromuserid string
	Url        string
	Ecode      string
	Ttype      uint8 // UNI, GRP
	Timet      int64
}
type InPut struct {
	Mtype    uint8
	Touserid string
	Url      string
	Ecode    string
	Ttype    uint8 // UNI, GRP
}

func (self Message) String() string {
	return string(self)
}
func ParseMessage(msg Message) *InPut {
	var val InPut
	temp := strings.Split(string(msg), "_")
	mtype := temp[0]
	switch mtype {
	case "emotion":
		val.Mtype = EMOTION
	case "picture":
		val.Mtype = PICTURE
	case "video":
		val.Mtype = VIDEO
	}
	val.Touserid = temp[1]
	if val.Mtype == EMOTION {
		val.Ecode = temp[2]
	} else {
		val.Url = temp[2]
	}

	ttype := temp[3]
	switch ttype {
	case "unicast":
		val.Ttype = UNI
	case "group":
		val.Ttype = GRP
	}
	return &val
}

func NewOutput(val InPut, fromuserid string) *OutPut {
	return &OutPut{val.Mtype, fromuserid, val.Url, val.Ecode, val.Ttype, time.Now().Unix()}
}
func (temp *OutPut) Bytes() []byte {
	val := make([]byte, 128)
	var n int
	buffer := bytes.NewBuffer(val)

	switch temp.Mtype {
	case EMOTION:
		n, _ = buffer.WriteString("emotion")
	case VIDEO:
		n, _ = buffer.WriteString("VIDEO")
	case PICTURE:
		n, _ = buffer.WriteString("PICTURE")
	}

	n, _ = buffer.WriteString("_")
	n, _ = buffer.WriteString(temp.Fromuserid)

	n, _ = buffer.WriteString("_")
	if temp.Mtype != EMOTION {
		n, _ = buffer.WriteString(temp.Url)
	} else {
		n, _ = buffer.WriteString(temp.Ecode)
	}

	n, _ = buffer.WriteString("_")
	if temp.Ttype == GRP {
		n, _ = buffer.WriteString("group")
	} else {
		n, _ = buffer.WriteString("unicast")
	}

	n, _ = buffer.WriteString("_")
	n, _ = buffer.WriteString(strconv.FormatInt(temp.Timet, 10))
	_ = n
	return buffer.Bytes()
}
