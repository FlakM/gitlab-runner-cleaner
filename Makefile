BINARY=gitlab-runner-cleaner.sh


VERSION = `git describe --tags`
BUILD = `date +%FT%T%z`

LDFLAGS = -ldflags "-w -s -X main.buildTag=$(VERSION) -X main.buildDate=$(BUILD)"

deps:
	$(call blue, "Downloading dependencies (using standalone binary)...")
	go get -t -v ./...
.PHONY: deps

build:
	$(call blue, "building...")
	go build $(LDFLAGS) -o $(BINARY)

test:
	$(call blue, "testing...")
	go test -v -race ./...
	go tool vet .

install:
	$(call blue, "installing...")
	go install ${LDFLAGS}

clean:
	$(call blue, "cleaning workspace...")
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install

define blue
	@tput setaf 6
	@echo $1
	@tput sgr0
endef
