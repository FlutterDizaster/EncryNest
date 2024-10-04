GENERATED_DIR = api/generated

PROTOC_GEN = protoc --go_out=. --go-grpc_out=.

PROTO_FILES := $(shell find . -name '*.proto')

BUILD_FOLDER = ./build

CLIENT_BINARY = client
SERVER_BINARY = server

CLIENT_PATH = ./cmd/client
SERVER_PATH = ./cmd/server

generate:
	$(PROTOC_GEN) $(PROTO_FILES)

clean-pb:
	rm -rf $(GENERATED_DIR)/*.pb.go

clean-build:
	rm -rf $(BUILD_FOLDER)

build-client:
	@echo "Building client..."
	go build -o $(BUILD_FOLDER)/$(CLIENT_BINARY) $(CLIENT_PATH)/main.go
	@echo "Client built successfully: $(CLIENT_BINARY)"

build-server:
	@echo "Building server..."
	go build -o $(BUILD_FOLDER)/$(SERVER_BINARY) $(SERVER_PATH)/main.go
	@echo "Server built successfully: $(SERVER_BINARY)"

client: build-client
server: build-server

.PHONY: clean-build clean-pb client server build-client build-server