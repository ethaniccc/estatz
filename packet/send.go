package packet

import (
	"bytes"
	"fmt"
	"net"
	"sync"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, MaxPacketSize))
	},
}

// Send writes a packet to a given address.
func Send(from *net.UDPConn, to *net.UDPAddr, header *PacketHeader, pk Packet) error {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()

	writer := protocol.NewWriter(buf, 0)
	header.Marshal(writer)
	pk.Encode(writer)

	pkSize := buf.Len()
	if pkSize > MaxPacketSize {
		// Here, we don't want to return the buffer to the pool because the cap of the buffer will be higher than
		// normal, and pools are the most effective when objects are of similar size.
		return fmt.Errorf("packet with ID %d with size %d exceeds MaxPacketSize (%d)", pk.ID(), pkSize, MaxPacketSize)
	}

	// Now that we know the cap of the buffer hasn't been exceeded, we can return this packet back to the buffer pool.
	defer bufferPool.Put(buf)
	if bytesWritten, err := from.WriteToUDP(buf.Bytes(), to); err != nil {
		return err
	} else if bytesWritten != pkSize {
		return fmt.Errorf("expected %d bytes to be written to connection, only wrote %d", pkSize, bytesWritten)
	} else {
		return nil
	}
}
