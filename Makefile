lint: 
    golangci-lint run ./... --config=./.golangci.yml

run:
	go run main.go

protogen:
	protoc \
        --go_out=proto_services \
        --go_opt=paths=source_relative \
	    --go-grpc_out=proto_services \
        --go-grpc_opt=paths=source_relative \
        handlers.proto

cov:
    go test -coverprofile c.out
    go tool cover -html c.out

bench:
     go test -run none -bench . -benchtime 3s -benchmem