.DEFAULT_GOAL := buildq

#
# Boring
#

fmt:
	go fmt ./...

lint: fmt
	golint ./...

vet: fmt
	go vet ./...

bootstrap:
	go mod download
	go get github.com/cosmtrek/air
	go get github.com/mitchellh/gox
	go install github.com/cosmtrek/air
	go install github.com/mitchellh/gox
	go generate -tags tools tools/tools.go

#
# Run, test, & bench
#

run:
	air

rundev:
	.\air.exe

test:
	go test --cover ./...

bench:
	go test --cover -bench . -benchmem ./...

#
# Build
#

buildq:
	go build -ldflags "-s -w" .

builddev:
	go build -ldflags "-s -w" -tags "music-api-dev" -o "bin/music-api-dev.exe" .

build: vet
	gox -osarch "linux/amd64 windows/amd64" \
	-gocmd go           \
	-ldflags "-s -w"    \
	-tags "music-api"    \
	-output "bin/music-api" .
