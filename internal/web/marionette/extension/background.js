// This file is loaded by the Quickbeam extension in the background,
// if the code below is uncommented it will redirect ALL websocket
// connections to the local quickbeam server, so they can be intercepted.

//function logURL(requestDetails) {
//  console.log(`Loading: ${requestDetails.url}`);
//  return {redirectUrl: `ws://localhost:18981?url=${encodeURIComponent(requestDetails.url)}`};
//};
//
//browser.webRequest.onBeforeRequest.addListener(logURL, {
//  urls: ["<all_urls>"],
//  types: ["websocket"]
//}, ["blocking"]);