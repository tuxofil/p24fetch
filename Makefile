.PHONY: all run test lint cover clean

NAME = p24fetch
COVER = .cover.out

all:
	go build -o $(NAME) ./cmd/main

run: all
	./$(NAME)

test:
	go test -v -coverprofile=$(COVER) -race ./...
	$(MAKE) lint

lint:
	golangci-lint run ./...

cover: test
	go tool cover -html=$(COVER)

clean:
	rm -rf -- $(NAME) $(COVER) run dedup/testdata
