package chat

import "io/ioutil"
import "net/http"

const NOTICEURL = "http://live.66boss.com/umeng/send?msgtype=chat&users="

func NoticeUMeng(user string, message string) {
	url := NOTICEURL + user + "&message=" + message
	res, err := http.Get(url)
	if err != nil {
		return
	}
	detail, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	_ = detail
}

//func main() {
//	NoticeUMeng("100000034", "HiGirl")
//}
