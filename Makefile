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
