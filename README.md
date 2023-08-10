[![builds.sr.ht status](https://builds.sr.ht/~michl/quickbeam.svg)](https://builds.sr.ht/~michl/quickbeam?)

Quickbeam runs web applications in a headless browser and provides APIs
to interact with them. Quickbeam APIs are machine readable but also
explorable and designed to be used by humans with minimal translation.

# Quickbeam-API

This sections refers to the API between user agents and quickbeam. This
API exposes the APIs for interacting with running web applications as
well as meta interactions like opening the web application in the first
place. It also exposes documentation functions about both meta and
applications level methods.

## Methods

### Meta

The `ping` method is simply answered by a `pong` reply.

The `version` method gets version information about quickbeam itself,
the browser and any activated modules.

### Navigation

The `open` method navigates to its single argument `url`. If the web
browser was not started yet, it is launched. If it is already running,
quickbeam forcefully leaves the current web page for the new `url`.

``` json
{ "jsonrpc": "2.0", "method": "open", "params": { "url": "http://example.com" }, "id": 1 }
```

The `close` method quits any running web browser and quickbeam itself.

``` json
{ "jsonrpc": "2.0", "method": "close", "id": 2 }
```

## Errors

Error categories brainstorm:

-   marionette.Page.Navigate: Reached error page (e.g. dnsNotFound)
-   Runtime error / panic

# Contact

info@infiniteaccess.eu
