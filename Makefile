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
	go get gotest.tools/gotestsum
	go install github.com/cosmtrek/air
	go install github.com/mitchellh/gox
	go install gotest.tools/gotestsum
	go generate -tags tools tools/tools.go

#
# Run, test, & bench
#

run:
	air

rundev:
	.\air.exe

test:
# go test --cover ./...
	gotestsum --format pkgname -- --cover ./...

bench:
	go test --cover -bench . -benchmem ./...

#
# Build
#

buildq:
	go build -ldflags "-s -w" .

builddev:
	go build -ldflags "-s -w" -tags "music-api-dev" -o "bin/music-api-dev.exe" .

builddocker:
	docker stop music-api
	docker rm music-api
	docker build --no-cache -t music-api .
	docker run -d --name music-api -p 80:80 music-api

build: vet
	gox -osarch "linux/amd64" \
	-gocmd go           \
	-ldflags "-s -w"    \
	-tags "music-api"    \
	-output "bin/music-api" .
