# Quickbeam extension

This extension contains JavaScript code that will be injected into all websites
that are visited by the Quickbeam Browser (Firefox).

## Loading the extension (Debug)
In the browser, navigate to [about:debugging#/runtime/this-firefox](the about:debugging page), under the tab "This Firefox", select "Load temporary addon", there select the [manifest.json](manifest.json) file.

## Loading the extension (Marionette)
The `Addon:Install` method can be used to install the extension. It accepts to arguments: `path` is the absolute path to the xpi file and `temporary` needs to be true as the extension is unsigned and can only be installed temporarily.

## Building the extension

Building the xpi file is straightforward as it is just a zip archive containing all relevant files. It is created and updated by running `zip -r -FS ../extension.xpi *` in the extension directory. The result is an unsigned webextension, so it can only be installed as temporary addons.
