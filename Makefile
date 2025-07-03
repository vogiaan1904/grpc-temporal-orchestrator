# Build variables
BUILD_DIR=./build

# Go build flags
GO_BUILD_FLAGS=-ldflags="-s -w"
APP_NAME=order-orchestrator

.PHONY: all build clean run-grpc run-http

run-server:
	@echo "Starting $(APP_NAME)..."
	go run cmd/server/main.go

protoc-all:
	$(MAKE) protoc PROTO=protos/proto/payment.proto OUT_DIR=protogen/golang/payment
	$(MAKE) protoc PROTO=protos/proto/order.proto OUT_DIR=protogen/golang/order
	$(MAKE) protoc PROTO=protos/proto/product.proto OUT_DIR=protogen/golang/product

protoc:
	protoc --go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	-I=protos/proto $(PROTO)