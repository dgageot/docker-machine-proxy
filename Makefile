BIN := docker-machine-proxy

export GO15VENDOREXPERIMENT = 1

.DEFAULT_GOAL := build

run: build
	./$(BIN)

build: $(BIN)

$(BIN): main.go machine/*.go proxy/*.go
	go build .

deps:
	godep save

clean:
	rm -f $(BIN)
