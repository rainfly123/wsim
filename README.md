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
