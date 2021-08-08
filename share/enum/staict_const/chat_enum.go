/**
 * @Author KYIMH
 * @Description
 * @Date 2021/8/2 16:49
 **/

package staict_const

type ChatMsg struct {
	ChatId  uint32 `bson:"chat_id"`  // unique id of this chat message
	Msg     []byte `bson:"msg"`      // message (may by encrypt)
	FromId  uint32 `bson:"from_id"`  // user id to send this message
	ToId    uint32 `bson:"to_id"`    // id of uer who can receive this message
	QueueId uint32 `bson:"queue_id"` // id of message queue
	MsgType uint8  `bson:"msg_type"` // message type(text, image, video...)
}
