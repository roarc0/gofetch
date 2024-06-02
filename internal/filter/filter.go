package filter

import (
	"errors"

	"github.com/roarc0/gct/internal/collector"
)

var (
	ErrUnknownMatchType = errors.New("unknown match type")
)

type Filter interface {
	Filter(dls []collector.Downloadable) ([]MatchedDownloadable, error)
}

type filterWithPartialMatches struct {
	matchers []Matcher
}

type MatchedDownloadable struct {
	collector.Downloadable
	Optional bool
}

// NewFilterWithOptionalMatches creates a new filter to determine which downloadables
// should be kept based on the matchers.
func NewFilterWithOptionalMatches(matchers []Matcher) Filter {
	return &filterWithPartialMatches{
		matchers: matchers,
	}
}

func (f *filterWithPartialMatches) Filter(in []collector.Downloadable) (out []MatchedDownloadable, err error) {
	for _, d := range in {
		match := true
		partialMatch := false

		for _, matcher := range f.matchers {
			m, err := matcher.Match(d)
			if err != nil {
				return nil, err
			}

			switch matcher.MatchType() {
			case MatchTypeRequired:
				if !m {
					match = false
				}
			case MatchTypeOptional:
				if m {
					partialMatch = true
				}
			case MatchTypeExclude:
				if m {
					match = false
				}
			default:
				return nil, ErrUnknownMatchType
			}
		}

		if match {
			out = append(out,
				MatchedDownloadable{
					Downloadable: d,
					Optional:     partialMatch,
				},
			)
		}
	}

	return out, nil
}

func FilterPartialMatchDownloadables(in []MatchedDownloadable) (out []collector.Downloadable) {
	for _, d := range in {
		if !d.Optional {
			out = append(out, d.Downloadable)
		}
	}

	return
}
