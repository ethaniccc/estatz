package packet

import (
	"bytes"
	"net"
	"sync"
)

var msgPool = sync.Pool{
	New: func() any {
		return &Message{
			buffer: bytes.NewBuffer(make([]byte, 0, MaxPacketSize)),
			sender: nil,
		}
	},
}

type Message struct {
	buffer *bytes.Buffer
	sender *net.UDPAddr
}

func NewMessage(buf []byte, sender *net.UDPAddr) *Message {
	msg := msgPool.Get().(*Message)
	msg.buffer.Write(buf)
	msg.sender = sender
	return msg
}

func (msg *Message) Buffer() *bytes.Buffer {
	return msg.buffer
}

func (msg *Message) Sender() *net.UDPAddr {
	return msg.sender
}

func (msg *Message) Dispose() {
	msg.sender = nil
	msg.buffer.Reset()
	msgPool.Put(msg)
}
