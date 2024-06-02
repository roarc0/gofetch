package collector

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMagnetCollectorCollect(t *testing.T) {
	body := `<a href="magnet:?xt=urn:btih:1234">magnet1</a><a href="magnet:?xt=urn:btih:5678">magnet2</a>`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}))
	defer srv.Close()

	tests := []struct {
		name string

		collector DownloadableCollector
		wantFn    func(ds []Downloadable, err error) error
	}{
		{
			name: "CollectOk",
			collector: func() DownloadableCollector {
				c, _ := NewMagnetCollector(srv.URL)
				return c
			}(),
			wantFn: func(dls []Downloadable, err error) error {
				if err != nil {
					return err
				}

				if len(dls) != 2 {
					return errors.New("expected 2 downloadables")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.collector
			got, err := m.Collect(context.Background())
			if wantErr := tt.wantFn(got, err); wantErr != nil {
				t.Errorf("MagnetCollect() error = %v", wantErr)
				return
			}
		})
	}
}
