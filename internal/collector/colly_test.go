package collector

import (
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/require"
)

func Test_httpClient(t *testing.T) {
	url := os.Getenv("MAGNETDL_URL")

	if url == "" {
		t.Skip("MAGNETDL_URL not set")
	}

	client := httpClient(&HttpConfig{})

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	req.Header.Set("User-Agent", defaultUserAgent)
	//req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)

	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	log.Println(string(body))

	require.Contains(t, string(body), "magnet:?")
}

func Test_colly(t *testing.T) {
	url := os.Getenv("MAGNETDL_URL")

	if url == "" {
		t.Skip("MAGNETDL_URL not set")
	}

	c := newColly(&HttpConfig{})
	count := 0
	c.OnHTML("a", func(e *colly.HTMLElement) {
		log.Println(e.Attr("href"))
		count++
	})

	err := c.Visit(url)
	require.NoError(t, err)
}
