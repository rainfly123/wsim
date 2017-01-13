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
	Mtype       uint8
	Fromuserid  string
	Url         string
	Ecode       string
	Ttype       uint8 // UNI, GRP
	Fromgroupid string
	Timet       int64
	extenion    string
}
type InPut struct {
	Mtype    uint8
	Touserid string //or fromgroupid
	Url      string
	Ecode    string
	Ttype    uint8 // UNI, GRP
	extenion string

	Fromuserid string
}

func (self Message) String() string {
	return string(self)
}
func ParseMessage(msg Message, from string) *InPut {
	var val InPut
	temp := strings.Split(string(msg), "_")
	if len(temp) != 5 {
		return nil
	}
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
	val.extenion = temp[4]
	val.Fromuserid = from
	return &val
}

func NewOutput(val *InPut) *OutPut {
	return &OutPut{val.Mtype, val.Fromuserid, val.Url, val.Ecode, val.Ttype, val.Touserid, time.Now().Unix(), val.extenion}
}
func (temp *OutPut) String() string {
	var val []string
	val = make([]string, 7)
	switch temp.Mtype {
	case EMOTION:
		val[0] = "emotion"
	case VIDEO:
		val[0] = "video"
	case PICTURE:
		val[0] = "picture"
	}
	val[1] = temp.Fromuserid
	if temp.Mtype != EMOTION {
		val[2] = temp.Url
	} else {
		val[2] = temp.Ecode
	}
	if temp.Ttype == GRP {
		val[3] = "group"
		val[4] = temp.Fromgroupid
	} else {
		val[3] = "unicast"
		val[4] = temp.Fromuserid
	}
	val[5] = strconv.FormatInt(temp.Timet, 10)
	val[6] = temp.extenion
	return strings.Join(val, "_")
}

func (temp *OutPut) Bytes() Message {
	val := make(Message, 256)
	var n int
	buffer := bytes.NewBuffer(val)

	switch temp.Mtype {
	case EMOTION:
		n, _ = buffer.WriteString("emotion")
	case VIDEO:
		n, _ = buffer.WriteString("video")
	case PICTURE:
		n, _ = buffer.WriteString("picture")
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
		n, _ = buffer.WriteString("_")
		n, _ = buffer.WriteString(temp.Fromgroupid)
	} else {
		n, _ = buffer.WriteString("unicast")
		n, _ = buffer.WriteString("_")
		n, _ = buffer.WriteString(temp.Fromuserid)
	}

	n, _ = buffer.WriteString("_")
	n, _ = buffer.WriteString(strconv.FormatInt(temp.Timet, 10))

	n, _ = buffer.WriteString("_")
	n, _ = buffer.WriteString(temp.extenion)
	_ = n
	return buffer.Bytes()
}

func WhetherLogin(msg Message) (bool, bool, string) {
	temp := string(msg)
	var action bool
	var login bool
	var userid string
	if strings.Contains(temp, "login") {
		action = true
		login = true
		userid = temp[strings.Index(temp, "_")+1:]
	} else if strings.Contains(temp, "logout") {
		action = true
		login = false
		userid = temp[strings.Index(temp, "_")+1:]
	}
	return action, login, userid
}
