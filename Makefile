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

bootstrap_ci:
# go mod download
	go get github.com/mitchellh/gox
	go get gotest.tools/gotestsum
	go install github.com/mitchellh/gox
	go install gotest.tools/gotestsum
	go generate -tags tools tools/tools.go

bootstrap: bootstrap_ci
	go get github.com/cosmtrek/air
	go install github.com/cosmtrek/air
	go generate -tags tools tools/tools.go

#
# Run, test, & bench
#

run:
	air

test:
	go clean -testcache
	gotestsum --format pkgname -- --cover -mod vendor ./...

bench:
	go test --cover -mod vendor -bench . -benchmem ./...

#
# Build
#

buildq:
	go build -ldflags "-s -w" -mod vendor .

builddev:
	go build -ldflags "-s -w" -mod vendor -tags "music-api-dev" -o "bin/music-api-dev.exe" .

builddocker:
	docker stop music-api
	docker rm music-api
	docker build --no-cache -t music-api .
#docker run -d --name music-api -p 80:80 music-api

build: vet
	gox -osarch "linux/amd64 windows/amd64" \
	-gocmd go           \
	-ldflags "-s -w"    \
	-tags "music-api"    \
	-output "bin/music-api" .
