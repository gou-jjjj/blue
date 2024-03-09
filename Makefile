bin:
	@rm -f ./bin/*
	@go build -o bin/ ./cmd/...
	@cp -r ./config.json bin/
	@mkdir -p bin/logs

commands:
	@go generate ./scripts/generate_commands.go

run:
	@go run ./cmd/main.go


req:run
	@go run ./api/go/main.go