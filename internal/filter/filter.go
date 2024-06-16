package filter

import (
	"errors"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/roarc0/gofetch/internal/collector"
)

var (
	ErrUnknownMatchType = errors.New("unknown match type")
)

type Filter struct {
	matchers []Matcher
}

type MatchedDownloadable struct {
	collector.Downloadable
	Optional bool
}

// NewFilter creates a new filter to determine which downloadables
// should be kept based on the matchers.
func NewFilter(matchers []Matcher) Filter {
	return Filter{
		matchers: matchers,
	}
}

func (f *Filter) Filter(in []collector.Downloadable) (out []MatchedDownloadable, err error) {
	for _, d := range in {
		match := true
		optionalMatch := false

		for _, matcher := range f.matchers {
			matchType, m, err := matcher.Match(d)
			if err != nil {
				return nil, err
			}

			switch matchType {
			case MatchTypeRequired:
				if !m {
					match = false
					goto next
				}
			case MatchTypeOptional:
				if m {
					optionalMatch = true
				}
			case MatchTypeSufficient:
				if m {
					match = true
					optionalMatch = false
					goto next
				}
			case MatchTypeExclude:
				if m {
					match = false
					goto next
				}
			default:
				return nil, ErrUnknownMatchType
			}
		}

	next:
		if match {
			out = append(out,
				MatchedDownloadable{
					Downloadable: d,
					Optional:     optionalMatch,
				},
			)
		}
	}

	return out, nil
}

func (f Filter) MarshalYAML() (any, error) {
	var fields struct {
		Matchers []MatcherWrapper
	}

	for _, m := range f.matchers {
		matcherType := reflect.TypeOf(m).String()
		switch matcherType {
		case reflect.TypeOf(&RegexMatcher{}).String():
			matcherType = "regex"
		default:
			return nil, errors.New("unknown matcher type")
		}

		fields.Matchers = append(
			fields.Matchers,
			MatcherWrapper{
				Type:    matcherType,
				Matcher: m,
			},
		)
	}

	return fields, nil
}

func (f *Filter) UnmarshalYAML(value *yaml.Node) error {
	var fields struct {
		Matchers []MatcherWrapper
	}

	if err := value.Decode(&fields); err != nil {
		return err
	}

	for _, mw := range fields.Matchers {
		f.matchers = append(f.matchers, mw.Matcher)
	}

	return nil
}
