package gofetch

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/roarc0/gofetch/internal/collector"
	"github.com/roarc0/gofetch/internal/config"
	"github.com/roarc0/gofetch/internal/filter"
	"github.com/roarc0/gofetch/internal/test/mocks"
)

var (
	//go:embed testdata/nyaa.html
	nyaaBody []byte
)

func TestGoFetch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(nyaaBody))
	}))

	cfg := config.Config{
		Sources: map[string]collector.Source{
			"nyaa": {
				Name: "nyaa",
				URIs: []string{srv.URL},
			},
		},
		Entries: map[string]filter.Entry{
			"test": {
				SourceName: "nyaa",
				Filter: filter.NewFilter(
					[]filter.Matcher{
						&filter.RegexMatcher{Regex: ".*Test.*"},
						&filter.RegexMatcher{Regex: ".*1080p.*", MatchType: filter.MatchTypeOptional},
						&filter.RegexMatcher{Regex: ".*480p.*", MatchType: filter.MatchTypeExclude},
					}),
			},
		},
	}

	mockMemory := mocks.MockMemory{}

	g, err := NewGoFetch(&cfg, &mockMemory)
	require.NoError(t, err)

	dls, err := g.Fetch()
	require.NoError(t, err)
	require.NotEmpty(t, dls)
}
