package collector

type Source struct {
	Name string
	URIs []string
}

func (s Source) Collector() (DownloadableCollector, error) {
	return Factory(s)
}
