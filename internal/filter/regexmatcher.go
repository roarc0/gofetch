package filter

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/roarc0/gofetch/internal/collector"
)

// RegexMatcher is a type that implements the Matcher interface to match a downloadable using a regex.
//
// NOTE: MatchType is used to determine what to do if the regex matches.
// We could use negative lookaheads in the regex, but go doesn't support them.
// This is because they can lead to denial of service attacks.
type RegexMatcher struct {
	Regex     string
	MatchType MatchType
}

func (m *RegexMatcher) Match(dl collector.Downloadable) (MatchType, bool, error) {
	matched, err := regexp.Match(m.Regex, []byte(dl.Name()))
	if err != nil {
		return MatchTypeInvalid, false, errors.Wrap(err, "failed to match regex")
	}

	return m.MatchType, matched, nil
}
