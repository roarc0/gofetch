package filter

import "gopkg.in/yaml.v3"

type MatchType int

const (
	MatchTypeRequired MatchType = iota
	MatchTypeOptional
	MatchTypeExclude
	MatchTypeSufficient
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
	case MatchTypeSufficient:
		return "sufficient"
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
	case "sufficient":
		*m = MatchTypeSufficient
	default:
		*m = MatchTypeInvalid
	}

	return nil
}

func (m MatchType) MarshalYAML() (any, error) {
	return m.String(), nil
}
