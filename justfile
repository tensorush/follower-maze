BIN_DIR := "bin"
BIN_FILE := "follower-maze"
BIN_PATH := BIN_DIR + "/" + BIN_FILE

all: verify build test

all-out: verify build-out test clean

verify:
    go mod verify

build:
    go build -v ./...

run:
    go run -v ./...

build-out:
    go build -o {{BIN_PATH}} cmd/follower-maze/main.go

run-out:
    ./{{BIN_PATH}}

test:
	go test -vet=off ./...

clean:
	go clean
	rm -rf {{BIN_DIR}}

docker-run:
    docker-compose up

docker-test:
    docker-compose run --rm follower-maze just test
