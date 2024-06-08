package collector

import (
	"context"
	_ "embed"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

//go:embed testdata/magnetdl.html
var magnetDlBody string

func TestMagnetDLMagnetCollectorCollect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(magnetDlBody))
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
				c, _ := NewMagnetDLMagnetCollector(srv.URL)
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
			c := tt.collector
			got, err := c.Collect(context.Background())
			if wantErr := tt.wantFn(got, err); wantErr != nil {
				t.Errorf("MagnetCollect() error = %v", wantErr)
				return
			}
		})
	}
}
