**Instant Message Based On WebSocket**

# 信令格式：
## **登录/登出：**

  1. login_userid  　用户登录
  2. logout_useid　用户退出，（为了偷懒可以不调用，服务端自动探测）

## **发送消息：**

  1. 发表情 emotion_touserid_表情编码_type_extension ; (touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  2. 发视频 video_touserid_URL_type_extension;(touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  3. 发图片 picture_touserid_URL_type_extension;(touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）

**extension 可以为任何字符　只要不含有_**
格式：
{“sender”:”张三”,”senderAvartar”:”www.baidu.com”,”conversation”:”老板六六聊天群”,”conversationAvartar”:”www.baidu.com”}sender:当前登录用户名,消息发送者senderAvartar:消息发送者头像conversation:私聊—》好友用户名，群聊—》群名conversationAvartar:私聊—》好友头像，群聊—》群头像

## **接收消息**

  1. 接收表情 emotion_fromuserid_表情编码_type_fromgroupid_timet_extension ; (fromuserid 是发送者用户ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊, timet是从1970年１月...到目前的秒数）
  2. 接收视频 video_fromuserid_URL_type_fromgroupid_timet_extension;(fromuserid 是发送者用户ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  3. 接收发图片 picture_fromuserid_URL_type_fromgroupid_timet_extension;(fromuserid 是发送者用户ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）

*如果是单聊，fromuserid = fromgroupid* 

## **查询(暂未实现，可能不需要用)**

 1. 查询成员，members_groupid
 2. 查询聊天记录, records_groupid/touserid_timet


## **心跳消息，服务器主动探测客户端，客户端收到后　不需要做任何处理**

 heartbeat_argument       

 argument 是变化的不定，后续扩展，只有heartbeat_　固定不变



## **上传图片**

http://live.66boss.com/upload/writev2?   Multipart-Form name="file"

## **上传视频**

http://live.66boss.com/upload/writev3?   Multipart-Form name="file"

## **聊天**
ws://live.66boss.com:6060/entry


##　群成员管理　
http://live.66boss.com:6060/creategrp?creator=xxx&members=abc,bcd,efg&name=x创建群

```
{
"code": 1,
"message": "OK",
"data": {
"groupid": "10",
"creator": "1000001653",
"name": "我的群",
"notice": "",
"snap": "http://live.66boss.com/emovideo/10.jpg",
"members": [
"1000001653",
"1000006331",
"1000006123",
"1000006340"
]
}
}
```

http://live.66boss.com:6060/addmember?groupid=xxx&members=abc,bcd,efg...增加群成员

{"code": 1,"message": "OK"}

http://live.66boss.com:6060/rmmember?groupid=xxx&members=abc,bcd,efg...踢出群成员

{"code": 1,"message": "OK"}

http://live.66boss.com:6060/editgrp?groupid=xxx&intro=xxx&notice=xxx&name=xx更改群资料

{"code": 1,"message": "OK"}

**查询我的群**
http://live.66boss.com:6060/querymygroups?userid=xxx 查询我的群

```
{
"code": 1,
"message": "OK",
"data": [
{
"groupid": "10",
"creator": "",
"name": "",
"notice": "",
"snap": "",
"members": [
"40",
"50",
"60",
"70",
"1000001653",
"1000006123",
"1000006331",
"1000006340"
]
},
{
"groupid": "11",
"creator": "",
"name": "",
"notice": "",
"snap": "",
"members": [
"1000001653",
"1000006123",
"1000006331",
"1000006340"
]
},
{
"groupid": "12",
"creator": "",
"name": "",
"notice": "",
"snap": "",
"members": [
"1000001653",
"1000006123",
"1000006331",
"1000006340"
]
},
{
"groupid": "13",
"creator": "",
"name": "",
"notice": "",
"snap": "",
"members": [
"50",
"60",
"1000001653",
"1000001901",
"1000006123",
"1000006331",
"1000006340"
]
}
]
}
{"code": 1,"message": "OK","data": ["5","8","3","","12","10","9","11","6","4","13","1","2","7"]}
```


