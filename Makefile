.PHONY: fmt lint test build clean hooks-install

fmt:
	mise run fmt

lint:
	mise run lint

test:
	mise run test

build:
	mise run build

clean:
	mise run clean

hooks-install:
	mise run hooks:install
