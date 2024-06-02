package filter

import "gopkg.in/yaml.v3"

type MatchType int

const (
	// ADD type Sufficient (doesn't need any optional matches to be kept)
	// Required will require at least one optional to match.

	MatchTypeRequired MatchType = iota
	MatchTypeOptional
	MatchTypeExclude
	MatchTypeInvalid
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
		return "invalid"
	}
}

func (m *MatchType) UnmarshalYAML(value *yaml.Node) error {
	var tmp string

	if err := value.Decode(&tmp); err != nil {
		return err
	}

	switch tmp {
	case "required":
		*m = MatchTypeRequired
	case "optional":
		*m = MatchTypeOptional
	case "exclude":
		*m = MatchTypeExclude
	default:
		*m = MatchTypeInvalid
	}

	return nil
}

func (m MatchType) MarshalYAML() (any, error) {
	return m.String(), nil
}
