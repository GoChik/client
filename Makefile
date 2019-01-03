VERSION = $(shell git describe --always)
GOFLAGS = -ldflags="-X main.Version=$(VERSION) -s -w"

.PHONY: default rpi_client gpio_client deploy help

default: help

rpi_client:
	GOOS=linux GOARCH=arm GOARM=6 go build -tags raspberrypi $(GOFLAGS)
	mkdir -p bin/client
	mv client bin/client/linux-arm

gpio_client:
	test -n "$(GOOS)" # GOOS
	test -n "$(GOARCH)" # GOARCH
	go build -tags gpio $(GOFLAGS)
	mkdir -p bin/client
	mv client bin/client/$(GOOS)-$(GOARCH)

fake_client:
	test -n "$(GOOS)" # GOOS
	test -n "$(GOARCH)" # GOARCH
	go build -tags fake $(GOFLAGS)
	mkdir -p bin/client
	mv client bin/client/$(GOOS)-$(GOARCH)

deploy:
	make rpi_client
	GOOS=linux GOARCH=mipsle make gpio_client
	GOOS=darwin GOARCH=amd64 make fake_client
	GOOS=linux GOARCH=amd64 make fake_client
	mkdir -p release/client
	rm -rf release/client/*
	@JFROG_CLI_OFFER_CONFIG=false jfrog bt dlv --user=rferrazz --key=$(BINTRAY_API_KEY) rferrazz/IO-Something/client/rolling release/
	go-selfupdate -o release/client bin/client $(VERSION)
	cd release && JFROG_CLI_OFFER_CONFIG=false jfrog bt u --user=rferrazz --key=$(BINTRAY_API_KEY) --override=true --flat=false --publish=true client/ rferrazz/IO-Something/client/rolling

clean:
	git clean -dfx

help:
	@echo "make [rpi_client gpio_client clean deploy]"
