package filter

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/roarc0/gofetch/internal/collector"
	"github.com/stretchr/testify/require"
)

func TestFilterFilter(t *testing.T) { //nolint:funlen
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
					collector.NewMagnet("test", "magnet:?xt=urn:btih:123"),
					collector.NewMagnet("abc", "magnet:?xt=urn:btih:123"),
				},
			},
			wantOut: []MatchedDownloadable{
				{
					Downloadable: collector.NewMagnet("test", "magnet:?xt=urn:btih:123"),
					Optional:     false,
				},
			},
		},
		{
			name: "FiltersIncludeExcludeOk",
			fields: fields{
				matchers: []Matcher{
					&RegexMatcher{Regex: "^t.*"},
					&RegexMatcher{Regex: "^.*1080p.*", MatchType: MatchTypeOptional},
					&RegexMatcher{Regex: "^.*480p.*", MatchType: MatchTypeExclude},
				},
			},
			args: args{
				in: []collector.Downloadable{
					collector.NewMagnet("test_1080p", "magnet:?xt=urn:btih:123"),
					collector.NewMagnet("test_720p", "magnet:?xt=urn:btih:123"),
					collector.NewMagnet("test_480p", "magnet:?xt=urn:btih:123"),
					collector.NewMagnet("abc", "magnet:?xt=urn:btih:123"),
				},
			},
			wantOut: []MatchedDownloadable{
				{
					Downloadable: collector.NewMagnet("test_1080p", "magnet:?xt=urn:btih:123"),
					Optional:     true,
				},
				{
					Downloadable: collector.NewMagnet("test_720p", "magnet:?xt=urn:btih:123"),
					Optional:     false,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Filter{
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

func TestFilter_MarshalUnmarshalYAML(t *testing.T) {
	type FilterWrapper struct {
		Filter Filter
	}

	for _, tv := range []struct {
		fw FilterWrapper
		s  string
	}{
		{
			FilterWrapper{
				Filter: NewFilter(
					[]Matcher{
						&RegexMatcher{Regex: "^t.*"},
						&RegexMatcher{Regex: "^v.*"},
					},
				),
			},
			"filter:\n    matchers:\n        - type: regex\n          matcher:\n            regex: ^t.*\n            matchtype: required\n        - type: regex\n          matcher:\n            regex: ^v.*\n            matchtype: required\n", //nolint:lll
		},
	} {
		t.Run(tv.s, func(t *testing.T) {
			b, err := yaml.Marshal(&tv.fw)
			require.NoError(t, err)
			require.Equal(t, tv.s, string(b))

			var fw FilterWrapper
			err = yaml.Unmarshal(b, &fw)
			require.NoError(t, err)
			require.Equal(t, tv.fw, fw)
		})
	}
}
