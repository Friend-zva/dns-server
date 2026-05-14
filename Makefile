.PHONY: install-tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: generate
generate:
	protoc -I proto/resolver \
		--go_out=proto/resolver \
	 	--go_opt=paths=source_relative \
		--go-grpc_out=proto/resolver \
		--go-grpc_opt=paths=source_relative \
		proto/resolver/resolver.proto

.PHONY: all
all: install-tools generate
