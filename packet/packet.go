package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	CurrentVersion = (0 << 24) | (0 << 16) | (0 << 8) | 1
)

// All packets sent to the EStatz server should be formatted as such:
// - Authentication JWT (??? bytes)
// - Client version (8 bytes)
// - Packet ID (4 bytes)
// - Packet data (??? bytes)
type Packet interface {
	ID() uint64
	Marshal(io protocol.IO, protoID uint64)
}

type PacketHeader struct {
	// JWT is the JWT token provided to prove that the client is authorized to send packets
	// to this server. This JWT can be authorized by the handler functions provided.
	JWT []byte
	// ClientVer represents the version the client is encoding/decoding packets from. This supports
	// backwards compatiability if a server update occurs.
	ClientVer uint64
	// PacketID is the ID of the packet to be decoded.
	PacketID uint64
}

func (h *PacketHeader) Marshal(io protocol.IO) {
	io.ByteSlice(&h.JWT)
	io.Uint64(&h.ClientVer)
	io.Uint64(&h.PacketID)
}
