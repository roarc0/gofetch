package collector

type Source struct {
	Type string
	URIs []string
}

func (s Source) Collector() (DownloadableCollector, error) {
	return Factory(s)
}
