
src_dir=proto/practical/v1
dest_dir=internal/practicalpb

all: clean resetdb test build

test: wire practicalpb
	go test -v ./...

build: wire practicalpb
	go build -o bin/server cmd/server/main.go

run: wire practicalpb
	APP_ENV=development \
	PORT=8080 \
	EVENT_STORE_CONNECTION_STRING="dbname=micro user=message_store password=postgres" \
	QUERY_CONNECTION_STRING=":6379" \
	go run cmd/server/main.go

wire:
	wire diff ./... | grep -q ^ && wire ./... || true

clean:
	rm -rf bin
	rm -f **/wire_gen.go
	@rm -f internal/application/identity/wire_gen.go
	@rm -f internal/application/viewing/wire_gen.go
	@rm -f internal/infrastructure/config/wire_gen.go
	@rm -f internal/presentation/rpc/wire_gen.go
	@rm -f internal/presentation/web/wire_gen.go
	@rm -rf $(dest_dir)

resetdb:
	psql -c "DROP DATABASE IF EXISTS micro"
	DATABASE_NAME=micro ~/Projects/message-db/database/install.sh
	redis-cli flushall

practicalpb: $(src_dir)/practical.proto
	mkdir -p $(dest_dir)
	protoc -I$(src_dir) \
		-I/usr/include/google/protobuf \
		--go_out $(dest_dir) \
		--go_opt paths=source_relative \
		--go-grpc_out $(dest_dir) \
		--go-grpc_opt paths=source_relative \
		$<
