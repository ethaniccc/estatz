package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

type Encodable interface {
	Encode(io *protocol.Writer)
}

type Decodable interface {
	Decode(io *protocol.Reader, version uint64)
}
