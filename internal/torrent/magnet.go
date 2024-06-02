package torrent

type Magnet struct {
	name string
	uri  string
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
