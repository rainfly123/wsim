Instant Message Based On WebSocket

##信令格式：
**登录/登出：**
  1. loginin_userid
  2. loginout_useid

**发送消息：**
  1. 发表情 emotion_touserid_表情编码_type ; (touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  2. 发视频 video_touserid_URL_type;(touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  3. 发图片 picture_touserid_URL_type;(touserid 是接收者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）


**接收消息：**
  1. 接收表情 emotion_fromuserid_表情编码_type_timet ; (fromuserid 是发送者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊, timet是从1970年１月...到目前的秒数）
  2. 接收视频 video_fromuserid_URL_type_timet;(fromuserid 是发送者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）
  3. 接收发图片 picture_fromuserid_URL_type_timet;(fromuserid 是发送者用户ID或者组ID，图片编码全网唯一，type 有二中:unicast, group 单聊，群聊）


**查询**
 1. 查询成员，members_groupid
 2. 查询聊天记录, records_groupid/touserid_timet
