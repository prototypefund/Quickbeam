all: build

XPI_BASE_PATH = internal/web/marionette
XPI_FILES = content.js background.js manifest.json icons/48.png icons/96.png

VERSION = 0.3

build:	internal/web/marionette/extension.xpi
	go build ./cmd/quickbeam

$(XPI_BASE_PATH)/extension.xpi: $(addprefix $(XPI_BASE_PATH)/extension/,$(XPI_FILES))
	cd internal/web/marionette/extension; zip -r -FS ../extension.xpi $(XPI_FILES)

extension: $(XPI_BASE_PATH)/extension.xpi

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
