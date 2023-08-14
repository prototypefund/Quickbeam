# Quickbeam extension

This extension contains JavaScript code that will be injected into all websites
that are visited by the Quickbeam Browser (Firefox).

## Loading the extension (Debug)
In the browser, navigate to [about:debugging#/runtime/this-firefox](the about:debugging page), under the tab "This Firefox", select "Load temporary addon", there select the [manifest.json](manifest.json) file.

## Loading the extension (Marionette)
The `Addon:Install` method can be used to install the extension. It accepts to arguments: `path` is the absolute path to the xpi file and `temporary` needs to be true as the extension is unsigned and can only be installed temporarily.

## Building the extension
TBD
