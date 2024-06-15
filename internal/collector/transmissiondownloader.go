package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"time"
)

// TransmissionDownloader is a type that implements the Downloader interface to open a downloadable using xdg-open.
type TransmissionDownloader struct {
	Downloadable
	Transmission *TransmissionConfig
	client       *http.Client
}

type TransmissionConfig struct {
	Host string
	Port string
	User string
	Pass string
	SSL  bool
}

type transmissionRequest struct {
	Method    string                `json:"method"`
	Arguments transmissionArguments `json:"arguments"`
}

type transmissionArguments struct {
	Paused   bool   `json:"paused"`
	Filename string `json:"filename"`
}

var (
	re = regexp.MustCompile("<code>X-Transmission-Session-Id: (.*)</code>")
)

// NewTransmissionDownloader creates a new TransmissionDownloader
func NewTransmissionDownloader(connection *TransmissionConfig) *TransmissionDownloader {
	return &TransmissionDownloader{
		Transmission: connection,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Download does an http post to a trasmission server to download the torrent file
func (d *TransmissionDownloader) Download(dl Downloadable) error {
	sessionID, err := d.getSessionID()
	if err != nil {
		return err
	}

	body, err := getBody(d.URI())
	if err != nil {
		return err
	}

	req, err := d.Transmission.newRequest()
	if err != nil {
		return err
	}
	req.Header.Set("X-Transmission-Session-Id", *sessionID)
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func getBody(uri string) ([]byte, error) {
	torrentAdd := transmissionRequest{
		Method: "torrent-add",
		Arguments: transmissionArguments{
			Paused:   false,
			Filename: uri,
		},
	}

	jsonStr, err := json.Marshal(torrentAdd)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

func (d *TransmissionDownloader) getSessionID() (*string, error) {
	req, err := d.Transmission.newRequest()
	if err != nil {
		return nil, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	match := re.FindStringSubmatch(string(body))
	if len(match) < 2 {
		return nil, fmt.Errorf("no session id found")
	}

	return &match[1], nil
}

func (tc *TransmissionConfig) newRequest() (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		tc.transmissionURL(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(tc.User, tc.Pass)

	return req, nil
}

func (tc *TransmissionConfig) transmissionURL() string {
	protocol := "http"
	if tc.SSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/transmission/rpc/",
		protocol,
		net.JoinHostPort(tc.Host, tc.Port),
	)
}
