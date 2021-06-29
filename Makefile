BIN= $(CURDIR)/bin
PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))

GO = go
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

PROTO_FILES_PATH=proto
PROTO_OUT=janusrpc

$(BIN):
	@mkdir -p $@

.PHONY: gprc
grpc: 
				protoc -I $(PROTO_FILES_PATH) --go_out=plugins=grpc:$(PROTO_OUT) $(PROTO_FILES_PATH)/*.proto

.PHONY: all
all: fmt $(BIN) ; $(info $(M) building executable) @ ## Build program binary
				$Q go build -o ./bin/janus -race ./cmd/main.go

.PHONY: fmt
fmt: ; $(info $(M) running gofmt...) @ ## Run gofmt on all source files
				$Q $(GO) fmt $(PKGS)

.PHONY: clean
clean: @rm -rf $(BIN)

.PHONY: run
run: 
				go run ./cmd/main.go