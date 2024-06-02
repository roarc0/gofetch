package filter

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

func TestMatchType_MarshalUnmarshalYAML(t *testing.T) {
	for _, tv := range []struct {
		mt MatchType
		s  string
	}{
		{MatchTypeRequired, "required"},
		{MatchTypeOptional, "optional"},
		{MatchTypeExclude, "exclude"},
		{MatchTypeInvalid, "invalid"},
	} {
		t.Run(tv.s, func(t *testing.T) {
			b, err := yaml.Marshal(&tv.mt)
			require.NoError(t, err)
			require.Contains(t, string(b), tv.s)

			var mt MatchType
			err = yaml.Unmarshal(b, &mt)
			require.NoError(t, err)
			require.Equal(t, tv.mt, mt)
		})
	}
}
