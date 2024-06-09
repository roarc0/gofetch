package filter

type Entry struct {
	SourceName string
	Disabled   bool
	Filter     Filter
}
