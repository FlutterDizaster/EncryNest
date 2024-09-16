GENERATED_DIR = api/generated

PROTOC_GEN = protoc --go_out=. --go-grpc_out=.

PROTO_FILES := $(shell find . -name '*.proto')

generate:
	$(PROTOC_GEN) $(PROTO_FILES)

clean:
	rm -rf $(GENERATED_DIR)/*.pb.go