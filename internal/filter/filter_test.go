package filter

import (
	"reflect"
	"testing"

	"github.com/roarc0/gct/internal/collector"
	"github.com/roarc0/gct/internal/torrent"
)

func TestFilterFilter(t *testing.T) {
	type fields struct {
		matchers []Matcher
	}
	type args struct {
		in []collector.Downloadable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantOut []MatchedDownloadable
		wantErr bool
	}{
		{
			name: "FilterOk",
			fields: fields{
				matchers: []Matcher{
					&RegexMatcher{Regex: "^t.*"},
				},
			},
			args: args{
				in: []collector.Downloadable{
					torrent.NewMagnet("test", "magnet:?xt=urn:btih:123"),
					torrent.NewMagnet("abc", "magnet:?xt=urn:btih:123"),
				},
			},
			wantOut: []MatchedDownloadable{
				{
					Downloadable: torrent.NewMagnet("test", "magnet:?xt=urn:btih:123"),
					Optional:     false,
				},
			},
		},
		{
			name: "FiltersIncludeExcludeOk",
			fields: fields{
				matchers: []Matcher{
					&RegexMatcher{Regex: "^t.*"},
					&RegexMatcher{Regex: "^.*1080p.*", MatchTypeVal: MatchTypeOptional},
					&RegexMatcher{Regex: "^.*480p.*", MatchTypeVal: MatchTypeExclude},
				},
			},
			args: args{
				in: []collector.Downloadable{
					torrent.NewMagnet("test_1080p", "magnet:?xt=urn:btih:123"),
					torrent.NewMagnet("test_720p", "magnet:?xt=urn:btih:123"),
					torrent.NewMagnet("test_480p", "magnet:?xt=urn:btih:123"),
					torrent.NewMagnet("abc", "magnet:?xt=urn:btih:123"),
				},
			},
			wantOut: []MatchedDownloadable{
				{
					Downloadable: torrent.NewMagnet("test_1080p", "magnet:?xt=urn:btih:123"),
					Optional:     true,
				},
				{
					Downloadable: torrent.NewMagnet("test_720p", "magnet:?xt=urn:btih:123"),
					Optional:     false,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &filterWithPartialMatches{
				matchers: tt.fields.matchers,
			}

			gotOut, err := f.Filter(tt.args.in)

			if (err != nil) != tt.wantErr {
				t.Errorf("filter.Filter() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("filter.Filter() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
