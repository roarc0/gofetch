package collector

func Factory(source Source) (DownloadableCollector, error) {
	switch source.Name {
	case "nyaa":
		return NewNyaaMagnetCollector(source.URIs[0])
	case "magnetdl":
		return NewMagnetDLMagnetCollector(source.URIs[0])
	default:
		return NewMagnetCollector(source.URIs[0])
	}
}
