package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMatcherWrapper_MarshalUnmarshalYAML(t *testing.T) {
	for _, tv := range []struct {
		mw MatcherWrapper
		s  string
	}{
		{
			mw: MatcherWrapper{
				Type:    "regex",
				Matcher: &RegexMatcher{Regex: "^t.*"},
			},
			s: "type: regex\nmatcher:\n    regex: ^t.*\n    matchtype: required\n",
		},
	} {
		t.Run(tv.s, func(t *testing.T) {
			b, err := yaml.Marshal(&tv.mw)
			require.NoError(t, err)
			require.Equal(t, tv.s, string(b))

			var mw MatcherWrapper
			err = yaml.Unmarshal(b, &mw)
			require.NoError(t, err)
			require.Equal(t, tv.mw, mw)
		})
	}
}
