package collector

import "time"

type Magnet struct {
	name string
	uri  string
	size uint64
	time time.Time
}

func NewMagnet(name, uri string) *Magnet {
	return &Magnet{
		name: name,
		uri:  uri,
	}
}

func (m Magnet) Name() string {
	return m.name
}

func (m Magnet) URI() string {
	return m.uri
}

func (m Magnet) String() string {
	return m.name + " " + m.uri
}

func (m Magnet) Size() uint64 {
	return m.size
}

func (m Magnet) Date() time.Time {
	return m.time
}
