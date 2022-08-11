build_wasnm:
	GOOS=js GOARCH=wasm go build -o server.wasm
caddy: build_wasnm
	GOROOT=$$(go env GOROOT) caddy run -watch