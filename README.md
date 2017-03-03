**Instant Message Based On WebSocket**

# 信令格式：
## **登录/登出：**

  1. login_userid  　用户登录
  2. logout_useid　用户退出，（为了偷懒可以不调用，服务端自动探测）

## **发送消息：**

  1. 发表情 emotion_touserid_表情编码_type_extension ; (touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  2. 发视频 video_touserid_URL_type_extension;(touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  3. 发图片 picture_touserid_URL_type_extension;(touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  4. 发文字 text_touserid_文字_type_extension ; (touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  5. 发音频 audio_touserid_URL_type_extension;(touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）

**extension 可以为任何字符　只要不含有_**
格式：
{“sender”:”张三”,”senderAvartar”:”www.baidu.com”,”conversation”:”老板六六聊天群”,”conversationAvartar”:”www.baidu.com”}sender:当前登录用户名,消息发送者senderAvartar:消息发送者头像conversation:私聊—》好友用户名，群聊—》群名conversationAvartar:私聊—》好友头像，群聊—》群头像

## **接收消息**

  1. 接收表情 emotion_fromuserid_表情编码_type_fromgroupid_timet_extension ; (fromuserid 是发送者用户ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊, timet是从1970年１月...到目前的秒数）
  2. 接收视频 video_fromuserid_URL_type_fromgroupid_timet_extension;(fromuserid 是发送者用户ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  3. 接收发图片 picture_fromuserid_URL_type_fromgroupid_timet_extension;(fromuserid 是发送者用户ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  3. 接收文字 text_fromuserid_文字_type_fromgroupid_timet_extension;(fromuserid 是发送者用户ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  2. 接收音频 audio_fromuserid_URL_type_fromgroupid_timet_extension;(fromuserid 是发送者用户ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）

*如果是单聊，fromuserid = fromgroupid* 

## **查询(暂未实现，可能不需要用)**

 1. 查询成员，members_groupid
 2. 查询聊天记录, records_groupid/touserid_timet


## **心跳消息，服务器主动探测客户端，客户端收到后　不需要做任何处理**

 heartbeat_argument       

 argument 是变化的不定，后续扩展，只有heartbeat_　固定不变


## **上传音频**
http://live.66boss.com/upload/writev1?   Multipart-Form name="file"


```
{
"code": 0,
"message": "OK",
"data": "http://livecdn.66boss.com/emopic//live/www/html/emopic/f85c0e8786bbc2d12db8d4f941f3dbd8-3.mp3"
#-3. 其中３是时长　单位秒
}
```

## **上传图片/音频**

http://live.66boss.com/upload/writev2?   Multipart-Form name="file"

```
{
"code": 0,
"message": "OK",
"data": "http://livecdn.66boss.com/emopic/1bada87cc98aa09399b38fb14628f54e-746x750.jpg"
}
```


## **上传视频**

```
http://live.66boss.com/upload/writev3?   Multipart-Form name="file"
{
"code": 0,
"message": "Succeeded",
"data": "http://livecdn.66boss.com/emovideo/76a0e06714b51d53bd93a64d1547c0b0-640x330.mp4"
}
```

## **聊天**
ws://live.66boss.com:6060/entry


##群管理　
创建群 http://live.66boss.com:6060/creategrp?creator=1000001653&members=1000001653,1000006331,1000006123,1000006340,1000001901&name=我的群 

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

增加群成员
http://live.66boss.com:6060/addmembers?groupid=10&members=50,60,70,40
{"code": 1,"message": "OK"}

踢出群成员
http://live.66boss.com:6060/delmembers?groupid=10&members=50,60    
{"code": 1,"message": "OK"}

更改群资料
http://live.66boss.com:6060/editgrp?groupid=10&notice=xxx&name=xx
{"code": 1,"message": "OK"}

查询我的群
http://live.66boss.com:6060/querymygrps?userid=xxx 查询我的群

```
{
"code": 1,
"message": "OK",
"data": [
{
"groupid": "26",
"creator": "100000050",
"name": "我的群",
"notice": "",
"snap": "http://live.66boss.com/emovideo/26_144fab26.jpg",
"members": [
{
"nickname": "18202093751",
"snap": "https://imgcdn.66boss.com/imagesu/avatar_temp/default.jpg"
},
{
"nickname": "18202093753",
"snap": "https://imgcdn.66boss.com/imagesu/avatar_temp/default.jpg"
},
{
"nickname": "段思琪",
"snap": "https://imgcdn.66boss.com/imagesu/avatar/20170216080727318380.jpg"
},
{
"nickname": "15323339887",
"snap": "https://imgcdn.66boss.com/imagesu/avatar_temp/default.jpg"
}
]
},
{
"groupid": "27",
"creator": "100000050",
"name": "我的群",
"notice": "",
"snap": "http://live.66boss.com/emovideo/27_a3c0cea6.jpg",
"members": [
{
"nickname": "18202093751",
"snap": "https://imgcdn.66boss.com/imagesu/avatar_temp/default.jpg"
},
{
"nickname": "18202093753",
"snap": "https://imgcdn.66boss.com/imagesu/avatar_temp/default.jpg"
},
{
"nickname": "段思琪",
"snap": "https://imgcdn.66boss.com/imagesu/avatar/20170216080727318380.jpg"
}
]
}]
}
```



## 查询群信息
http://live.66boss.com:6060/grpinfo?groupid=xxx 查询群信息

## 所有ＡＰＩ　可以用
http://live.66boss.com/wsim/??? 方式请求
