package collector

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const magnetDlBody = `
<tr>
   <td class="m"><a href="magnet:?xt=urn:btih:3ead8471d0d248d4bba76349bc2d367bfc5284d7&amp;dn=Test+%5Beztv%5D&amp;tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337" title="Direct Download" rel="nofollow"><img src="/img/m.gif" alt="Magnet Link" width="14" height="17" /></a></td>
   <td class="n"><a href="/file/5988604/test.2024.s01e03.1080p.hevc.x265-aaa/" title="Test.2024.S01E03.1080p.HEVC.x265-AAA">Test.2024.S01E03.1080p.HEVC.x265-AAA-<b>aaa</b></a></td>
   <td>2 months</td>
   <td class="t5">TV</td>
   <td>1</td>
   <td>754.08 MB</td>
   <td class="s">496</td>
   <td class="l">82</td>
</tr>
`

func TestMagnetDLMagnetCollectorCollect(t *testing.T) {
	t.Skip()
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

				if len(dls) != 1 {
					return errors.New("expected 1 downloadables")
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
