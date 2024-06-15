package collector

import (
	"testing"
)

func TestTransmissionDownloader_Open(t *testing.T) {
	t.Skip()

	magnet := Magnet{
		uri: "magnet:?xt=urn:btih:88594aaacbde40ef3e2510c47374ec0aa396c08e&dn=Big+Buck+Bunny+1080p+30fps&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451&tr=udp%3A%2F%2Fopen.stealth.si%3A80&tr=udp%3A%2F%2Ftracker.moeking.me%3A6969&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969&tr=udp%3A%2F%2Fopen.tracker.cl%3A1337",
	}
	td := NewTransmissionDownloader(
		&TransmissionConfig{
			Host: "localhost",
			Port: "9091",
			User: "transmission",
			Pass: "",
		},
	)

	td.Download(magnet)
}
