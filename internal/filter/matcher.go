package filter

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/roarc0/gofetch/internal/collector"
	"gopkg.in/yaml.v3"
)

// Matcher is an interface that defines a method to match a downloadable to see if it should be used or not.
type Matcher interface {
	Match(dl collector.Downloadable) (MatchType, bool, error)
}

type MatcherWrapper struct {
	Matcher
	Type string
}

func (m MatcherWrapper) MarshalYAML() (any, error) {
	matcherType := reflect.TypeOf(m.Matcher).String()
	switch matcherType {
	case reflect.TypeOf(&RegexMatcher{}).String():
		matcherType = "regex"
	default:
		return nil, errors.New("unknown matcher type")
	}

	return struct {
		Type    string
		Matcher Matcher
	}{
		Type:    matcherType,
		Matcher: m.Matcher,
	}, nil
}

func (m *MatcherWrapper) UnmarshalYAML(node *yaml.Node) error {
	var tmp struct {
		Type    string
		Matcher map[string]any
	}

	if err := node.Decode(&tmp); err != nil {
		return err
	}

	m.Type = tmp.Type

	switch m.Type {
	case "regex":
		m.Matcher = &RegexMatcher{}

		b, err := yaml.Marshal(tmp.Matcher)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(b, m.Matcher)

	default:
		return errors.Errorf("unknown matcher type %s", m.Type)
	}
}
