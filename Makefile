
all:
	go run ./cmd/gofetch-cli -mode manual

debug:
	dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient ./cmd/gofetch-cli -- -mode manual
