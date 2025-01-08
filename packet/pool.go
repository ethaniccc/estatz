package packet

type PacketFunc func() Packet

var pkPool = make(map[uint64]PacketFunc)

func Register(pk PacketFunc) {
	pkPool[pk().ID()] = pk
}

func Find(id uint64) (Packet, bool) {
	pkFunc, ok := pkPool[id]
	if !ok {
		return nil, false
	}
	return pkFunc(), true
}
