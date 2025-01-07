package packet

import (
	"bytes"
	"net"
	"sync"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

var msgPool = sync.Pool{
	New: func() any {
		return &Message{
			buffer: bytes.NewBuffer(make([]byte, 0, 1492)),
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

func (msg *Message) Sender() *net.UDPAddr {
	return msg.sender
}

func (msg *Message) Decode() (*PacketHeader, Packet, bool) {
	reader := protocol.NewReader(msg.buffer, 0, false)
	header := &PacketHeader{}
	header.Marshal(reader)
	pk, packetFound := Find(header.PacketID)
	pk.Marshal(reader, header.ClientVer)
	return header, pk, packetFound
}

func (msg *Message) Dispose() {
	msg.sender = nil
	msg.buffer.Reset()
	msgPool.Put(msg)
}
