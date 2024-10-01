GOOS=linux
GOARCH=amd64

NAME=apt-remove-version-$(GOOS)-$(GOARCH)

build:
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a \
		-ldflags "-w -s" \
		-trimpath \
		-o $(NAME)

build-arm: GOARCH=arm
build-arm: NAME=apt-remove-version-$(GOOS)-$(GOARCH)
build-arm: build

all:
	$(MAKE) build
	$(MAKE) build-arm

.DEFAULT_GOAL := all
