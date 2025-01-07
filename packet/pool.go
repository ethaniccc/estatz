package packet

type pkFunc func() Packet

var pkPool map[uint64]pkFunc

func Register(pk pkFunc) {
	pkPool[pk().ID()] = pk
}

func Find(id uint64) (Packet, bool) {
	pkFunc, ok := pkPool[id]
	if !ok {
		return nil, false
	}
	return pkFunc(), true
}
