[![builds.sr.ht status](https://builds.sr.ht/~michl/quickbeam.svg)](https://builds.sr.ht/~michl/quickbeam?)

Quickbeam runs web applications in a headless browser and provides APIs to
interact with them. Quickbeam APIs are machine readable but also explorable and
designed to be used by humans with minimal translation. The primary intention of
this tool is the enable the creation of accessibility focused clients to widely
used web applications that lack both accessibility and public APIs.

# Installation

To install quickbeam, build the Go package `cmd/quickbeam` in this repo and copy
the resulting `quickbeam` binary to a place from where it can be invoked. Go
1.18 is required. The most convenient way to do this is to run `make build`. If
you do not have a go compiler but can run docker or podman, you can also use
`./build.sh`. This script builds a Docker image containing everything to build
quickbeam and runs `make build` inside it.

To run the test suite, there are `check` and `check-short` targets in the
Makefile, executing all existing tests. Some test spawn a full browser instance
which takes time and other resources. These tests are skipped with
`check-short`.

The binaries produced will most likely dynamically depend on libc. So if you use
NixOS, Alpine Linux or another Linux distribution where libc is somehow special,
the binaries might not be as portable as perhaps expected. If you want to
distribute them, the Docker / podman way is a quick way to produce binaries that
work for Debian and most mainstream Linux distributions.

# Development

A development environment is easy to setup. The repository contains a Nix flake
which provides a development shell for the project. There is nothing special in
it though - only go, go tooling and web browsers. It is therefore easily
recreated without Nix.

There is a small script `./quicbeam-run.sh` which launches a quickbeam instance
via `go run` and duplicates the incoming and outgoing JSONRPC traffic to the
files `protocol-in.log` and `protocol-out.log` respectively. This is handy for
debugging a client or working on the API in quickbeam.

# Usage

Quickbeam exposes a JSON-RPC-API via standard input and standard output. It accepts a single command line flag `--headless` (default) or `--no-headless` which determines how the browser instance is launched. Errors and debugging information is logged to standard error. It is meant to be called by a client program that consumes the API and offers an interface to the user.

# Quickbeam-API

This sections refers to the API between user agents and quickbeam. This API
exposes the APIs for interacting with running web applications as well as meta
interactions like opening the web application in the first place. It also
exposes documentation functions about both meta and applications level methods.

## Meta

The `ping` method is simply answered by a `pong` reply.

The `version` method gets version information about quickbeam itself, the
browser and any activated modules.

The `wait` method waits a number of seconds (given by the `for` argument), then answers. It is used for testing clients.

## Navigation

The `open` method navigates to its single argument `url`. It forcefully
leaves the current web page for the new `url`.

``` json
{ "jsonrpc": "2.0", "method": "open", "params": { "url": "http://example.com" }, "id": 1 }
```

## Applications

Application level APIs consist of actions, objects, collections, documents and forms. Some methods to operate on these are:

- `call` with parameters `action` and `params`. The later is an object containing the parameters needed by the action. The result value it returns also depends on the action called.
- `fetch` with parameter `collection` returns all objects in the specified collection.
- `subscribe`with parameter `collection` registers the collection as subscribed. It returns the `id` of the subscription. Whenever the collection changes in the future, a JSON-RPC notification will be send to the client:

``` json
{ "jsonrpc": "2.0",
  "method": "collection_change",
  "params": {
    "collection": "some_collection",
    "id": "theid",
  }
}
```
- `state` returns which state the web application is in. The exact value depends on the application.

So far, the only application API that has been implement is `bbb`, which provides an interface to BigBlueButton conferences.

# Contact

[info@infiniteaccess.eu](mailto:info@infiniteaccess.eu)

# Acknowledgements

Quickbeam development was funded from March 2023 until September 2023 by Prototype Fund and the German Federal Ministry
of Education and Research.

![logos of the "Bundesministerium für Bildung und Forschung", Prototype Fund and Open Knowledge Foundation Deutschland](doc/pf_funding_logos.svg)
