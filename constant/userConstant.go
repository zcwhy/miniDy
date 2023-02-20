package constant

const (
	MAX_USERNAME_LEN = 32
	MAX_PASSWORD_LEN = 32
)

// MAX_VIDEO_NUMBER 控制单次返回的最大视频数量
const MAX_VIDEO_NUMBER = 30

// MAX_MESSAGE_NUMBER 控制单次返回的最大消息记录数量
const MAX_MESSAGE_NUMBER = 30

var (
	// FROM_MESSAGE 表示消息的类型是当前用户接收的消息
	FROM_MESSAGE = int64(0)
	// TO_MESSAGE 表示消息的类型是当前用户发送的消息
	TO_MESSAGE = int64(1)
)
