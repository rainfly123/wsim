package main

import "fmt"
import "io/ioutil"
import "net/http"
import "strings"
import "net/url"
import "encoding/base64"
import "encoding/json"

const NOTICEURL = "http://live.66boss.com/umeng/send?"

func NoticeUMeng(user string, message string) {
	query := url.Values{"users": {user}, "msgtype": {"chat"}, "message": {message}}
	url := query.Encode()
	res, err := http.Get(NOTICEURL + url)
	if err != nil {
		return
	}
	detail, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	_ = detail
}

func ParseString(msg string) (string, string) {
	i := strings.LastIndex(msg, "_")
	if i < 0 {
		return "", ""
	}
	i += 1
	uDec, err := base64.StdEncoding.DecodeString(msg[i:])
	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	all := make(map[string]interface{})
	json.Unmarshal(uDec, &all)
	result := all["conversation"]
	conversation := result.(string)
	var detail string
	switch {
	case strings.Contains(msg, "text"):
		{
			i := strings.Index(msg, "_")
			i += 1
			msg = msg[i:]
			j := strings.Index(msg, "_")
			j += 1
			msg = msg[j:]
			k := strings.Index(msg, "_")
			detail = ":" + msg[:k]
		}
	case strings.Contains(msg, "emotion"):
		{
			detail = ":[表情]"
		}
	case strings.Contains(msg, "picture"):
		{
			detail = ":[图片]"
		}
	case strings.Contains(msg, "video"):
		{
			detail = ":[视频]"
		}
	case strings.Contains(msg, "audio"):
		{
			detail = ":[语音]"
		}
	}

	return conversation, detail
}

/*
func main() {
	message, detail := ParseString("text_100000048_\xe5\x95\xa6\xe5\x95\xa6\xe5\x95\xa6\xe5\xbe\xb7\xe6\x81\xb6\xe9\xad\x94_unicast_100000048_1490938556_eyJzZW5kZXIiOiLmn6DmqqwiLCJzZW5kZXJJRCI6IjEwMDAwMDA0OCIsInNlbmRlckF2YXJ0YXIiOiJodHRwczpcL1wvaW1nY2RuLjY2Ym9zcy5jb21cL2ltYWdlc3VcL2F2YXRhclwvMjAxNzAzMjgxMDIzMDY5MjA2NjIuanBnIiwiY29udmVyc2F0aW9uIjoi5p+g5qqsIiwiY29udmVyc2F0aW9uQXZhcnRhciI6Imh0dHBzOlwvXC9pbWdjZG4uNjZib3NzLmNvbVwvaW1hZ2VzdVwvYXZhdGFyXC8yMDE3MDMyODEwMjMwNjkyMDY2Mi5qcGcifQ==")
	fmt.Println(detail)
	NoticeUMeng("100000073", message+detail)

}
*/
