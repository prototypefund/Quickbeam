# Quickbeam extension

This extension contains JavaScript code that will be injected into all websites
that are visited by the Quickbeam Browser (Firefox).

## Loading the extension (Debug)
In the browser, navigate to [about:debugging#/runtime/this-firefox](the about:debugging page), under the tab "This Firefox", select "Load temporary addon", there select the [manifest.json](manifest.json) file.

## Loading the extension (Marionette)
The `addon:install` method can be pointed to the (absolute) path of the [manifest.json](manifest.json).

## Building the extension
TBD
