protogen:
	protoc \
    --go_out=proto_services \
    --go_opt=paths=source_relative \
	--go-grpc_out=proto_services \
    --go-grpc_opt=paths=source_relative \
    handlers.proto

run:
    go run main.go