package websocketserver

// MessageType 定义消息类型的常量
const (
	MsgTypeEcho       = "echo"
	MsgTypeBroadcast  = "broadcast"
	UPDATE_CHARACTERS = "update_characters"
)

// Message 定义消息结构
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
