.PHONY: build run clean dev

build:
	cd apps/server-v2 && go build -o bin/server-v2 .

run: build
	cd apps/server-v2 && ./bin/server-v2

dev:
	cd apps/server-v2 && go run .

clean:
	rm -rf apps/server-v2/bin/ apps/server-v2/data/
