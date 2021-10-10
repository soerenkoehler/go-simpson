.PHONY: build

build:
	go get -d github.com/soerenkoehler/simpson@main
	go run github.com/soerenkoehler/simpson ${SIMPSON_CMDLINE}
	go mod tidy
