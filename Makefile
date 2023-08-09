all: build

VERSION = 0.2

build:
	go build ./cmd/quickbeam

package: quickbeam-$(VERSION).gz
quickbeam-$(VERSION).gz: build
	gzip -k -f quickbeam
	mv quickbeam.gz quickbeam-$(VERSION).gz

clean:
	rm -f quickbeam quickbeam.gz

check:
	go test -v ./internal/api
	go test -v ./internal/bbb
	go test -v ./internal/web/marionette
	#go test -v ./internal/web/rod

check-short:
	go test -v -test.short ./internal/api
	go test -v -test.short ./internal/bbb
	go test -v -test.short ./internal/web/marionette
	#go test -v -test.short ./internal/web/rod
