BIN := docker-machine-proxy

export GO15VENDOREXPERIMENT = 1

.DEFAULT_GOAL := build

run: build
	./$(BIN)

build: $(BIN)

$(BIN): main.go machine/*.go proxy/*.go
	GOOS=windows GOARG=amd64 go build -ldflags "-extldflags -static" .

deps:
	godep save

clean:
	rm -f $(BIN)
