package chat

import "time"

//import "log"
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
	Mtype      uint8
	Touserid   string
	Url        string
	Ecode      string
	Ttype      uint8 // UNI, GRP
	Fromuserid string
}

func (self Message) String() string {
	return string(self)
}
func ParseMessage(msg Message, from string) *InPut {
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
	default:
		return nil
	}
	val.Touserid = temp[1]
	if val.Mtype == EMOTION {
		val.Ecode = temp[2]
	} else {
		val.Url = temp[2]
	}

	ttype := temp[3]
	if strings.Contains(ttype, "unicast") {
		val.Ttype = UNI
	} else {
		val.Ttype = GRP
	}
	val.Fromuserid = from
	return &val
}

func NewOutput(val *InPut) *OutPut {
	return &OutPut{val.Mtype, val.Fromuserid, val.Url, val.Ecode, val.Ttype, time.Now().Unix()}
}
func (temp *OutPut) Bytes() Message {
	val := make(Message, 128)
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

func WhetherLogin(msg Message) (bool, bool, string) {
	temp := string(msg)
	var action bool
	var login bool
	var userid string
	if strings.Contains(temp, "loginin") {
		action = true
		login = true
		userid = temp[strings.Index(temp, "_")+1:]
	} else if strings.Contains(temp, "loginout") {
		action = true
		login = false
		userid = temp[strings.Index(temp, "_")+1:]
	}
	return action, login, userid
}
