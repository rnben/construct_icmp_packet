.PHONY: all build run

SERVER_BIN=icmp_test

all: build run
build:
	@go build -o test/bin/"${SERVER_BIN}" test/main.go
run:
	@sudo ./test/bin/"${SERVER_BIN}"