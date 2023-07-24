all: build

build:
	go build ./cmd/quickbeam

package: build
	gzip -f quickbeam

clean:
	rm -f quickbeam quickbeam.gz
