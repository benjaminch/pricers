MAIN_VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1")
VERSION:=${MAIN_VERSION}\#$(shell git log -n 1 --pretty=format:"%h")
PACKAGES:=$(shell go list ./... | sed -n '1!p' | grep -v /vendor/)
LDFLAGS:=-ldflags "-X github.com/benjaminch/pricers/app.Version=${VERSION}"

default: run

## depends: Installs dependencies running `dep ensure` internally
depends:
	rm -rf vendor
	go get -v ./...
	go build -v ./...
	go mod tidy

## upgrade: Upgrade the repository
upgrade:
	go get -u

## test: Runs tests running `validate.sh` script internally along with tests coverage
test:
	echo "mode: count" > coverage-all.out
	./scripts/tests/validate.sh
	# $(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

## cover: Runs tests coverage and output it in `coverage-all.out`
cover: test
	go tool cover -html=coverage-all.out

## run: Runs server in test env
run: build
	go run ${LDFLAGS} cmd/server/gondolier.go --env=debug

build: clean depends
	go build ${LDFLAGS} -a ./...

## clean: Cleans executable and tests coverage report files
clean:
	go clean
	rm -rf coverage.out coverage-all.out

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo