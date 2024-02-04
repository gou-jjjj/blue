bin:
	go build -o bin/ ./cmd/...
	cp -r ./config.json bin/
	mkdir -p bin/logs

commands:
	go generate ./scripts/generate_commands.go

