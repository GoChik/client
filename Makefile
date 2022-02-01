VERSION = $(shell git describe --always)
GOFLAGS = -ldflags="-X main.Version=$(VERSION) -s -w"

.PHONY: default dependencies rpi_client gpio_client deploy bintray-deploy help

default: help

dependencies:
	go get -u github.com/rferrazz/go-selfupdate

raspberrypi_client:
	test -n "$(GOOS)" # GOOS
	test -n "$(GOARCH)" # GOARCH
	go build -tags platform_raspberrypi $(GOFLAGS)
	mkdir -p bin/client
	mv client bin/client/$(GOOS)-$(GOARCH)

unipi_client:
	test -n "$(GOOS)" # GOOS
	test -n "$(GOARCH)" # GOARCH
	go build -tags platform_unipi_neuron $(GOFLAGS)
	mkdir -p bin/client
	mv client bin/client/$(GOOS)-$(GOARCH)

soft_client:
	test -n "$(GOOS)" # GOOS
	test -n "$(GOARCH)" # GOARCH
	go build -tags platform_soft $(GOFLAGS)
	mkdir -p bin/client
	mv client bin/client/$(GOOS)-$(GOARCH)

deploy:
	GOOS=linux GOARCH=arm GOARM=6 make raspberrypi_client
	GOOS=linux GOARCH=mipsle make raspberrypi_client
	GOOS=darwin GOARCH=amd64 make soft_client
	GOOS=linux GOARCH=amd64 make soft_client

bintray-deploy:
	deploy
	mkdir -p release/client
	rm -rf release/client/*
	@JFROG_CLI_OFFER_CONFIG=false jfrog bt dlv --user=rferrazz --key=$(BINTRAY_API_KEY) rferrazz/IO-Something/client/rolling release/
	go-selfupdate -o release/client bin/client $(VERSION)
	cd release && JFROG_CLI_OFFER_CONFIG=false jfrog bt u --user=rferrazz --key=$(BINTRAY_API_KEY) --override=true --flat=false --publish=true client/ rferrazz/IO-Something/client/rolling

clean:
	git clean -dfx

help:
	@echo "make [ dependencies rpi_client gpio_client deploy bintray-deploy clean ]"
