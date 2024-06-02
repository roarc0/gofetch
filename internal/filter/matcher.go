package filter

import (
	"regexp"

	"github.com/pkg/errors"

	"github.com/roarc0/gct/internal/collector"
)

type MatchType int

const (
	MatchTypeRequired MatchType = iota
	MatchTypeOptional
	MatchTypeExclude
)

func (m MatchType) String() string {
	switch m {
	case MatchTypeRequired:
		return "required"
	case MatchTypeOptional:
		return "optional"
	case MatchTypeExclude:
		return "exclude"
	default:
		return "unknown"
	}
}

// Matcher is an interface that defines a method to match a downloadable to see if it should be used or not.
type Matcher interface {
	Match(dl collector.Downloadable) (bool, error)
	MatchType() MatchType
}

// RegexMatcher is a type that implements the Matcher interface to match a downloadable using a regex.
//
// NOTE: Exclude is used to determine if the matcher should be used to exclude or include downloadables.
// We could use negative lookaheads in the regex, but go doesn't support them.
// This is because they can lead to denial of service attacks.
type RegexMatcher struct {
	Regex        string    `json:"regex"`
	MatchTypeVal MatchType `json:"match_type"`
}

func (m *RegexMatcher) Match(dl collector.Downloadable) (bool, error) {
	matched, err := regexp.Match(m.Regex, []byte(dl.Name()))
	if err != nil {
		return false, errors.Wrap(err, "failed to match regex")
	}

	return matched, nil
}

func (m *RegexMatcher) MatchType() MatchType {
	return m.MatchTypeVal
}
