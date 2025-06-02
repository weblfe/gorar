VERSION := 1.0.0

.PHONY: all
# build for windows, mac, linux
all: build-win build-mac build-linux

.PHONY: build-mac
# build  for mac
build-mac:
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=${VERSION}" -o bin/mac/gorar

.PHONY: build-linux
# build  for linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o bin/linux/gorar

.PHONY: build-win
# build  for windows
build-win:
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=${VERSION}" -o bin/win/gorar.exe

.PHONY: clean
# clean
clean:
	@rm -rf bin

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help