package collector

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransmissionDownloader_Open(t *testing.T) {
	user := os.Getenv("TRANSMISSION_USER")
	pass := os.Getenv("TRANSMISSION_PASS")

	if user == "" || pass == "" {
		t.Skip("TRANSMISSION_USER and TRANSMISSION_PASS are not set")
	}

	magnet := Magnet{
		uri: "magnet:?xt=urn:btih:88594aaacbde40ef3e2510c47374ec0aa396c08e&dn=Big+Buck+Bunny+1080p+30fps&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451&tr=udp%3A%2F%2Fopen.stealth.si%3A80&tr=udp%3A%2F%2Ftracker.moeking.me%3A6969&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969&tr=udp%3A%2F%2Fopen.tracker.cl%3A1337",
	}

	td := NewTransmissionDownloader(
		&TransmissionConfig{
			Host: "localhost",
			Port: "9091",
			User: user,
			Pass: pass,
		},
	)

	err := td.Download(magnet)
	require.NoError(t, err)
}
