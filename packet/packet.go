package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

var MaxPacketSize int = 1200

type Packet interface {
	Encodable
	Decodable

	ID() uint64
}

type PacketHeader struct {
	// Passphrase is a token provided to prove that the client is authorized to send packets
	// to this server. This should be handled by the specified handlers.
	Passphrase []byte
	// Version represents the version the sender is encoding/decoding packets from. This supports
	// backwards compatiability if a server update occurs, but the client is still running on an
	// older version.
	Version uint64
	// PacketID is the ID of the packet to be decoded.
	PacketID uint64
}

func (h *PacketHeader) Marshal(io protocol.IO) {
	io.ByteSlice(&h.Passphrase)
	io.Uint64(&h.Version)
	io.Uint64(&h.PacketID)
}
