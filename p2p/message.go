package p2p

const (
	IncomingMessage = 0x1
	IncomingStream  = 0x2
)

// RPC 保存通过网络中两个节点之间的
// 网络中两个节点之间的每种传输方式发送的任意数据。
type RPC struct {
	From    string
	Payload []byte
	Stream  bool
}