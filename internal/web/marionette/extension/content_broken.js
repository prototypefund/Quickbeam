// NOTE: This file is unused, it contains my debugging journey
// and is kept for reference.

// The original intent was to wrap the websocket constructor,
// so the extension could intercept all websocket connections,
// and add custom (onMessage) handlers.
// This turned out to be impossible, because the wrappedJSObject
// is not available in the context of the website, because of
// the "xRayVision" security feature of Firefox.

// We have the following options to pursue in the future:
//   1. Redirect all connections to websocket server. Intercept on go side.
//   2. Rewrite websocket constructer (that's what was tried in this file,
//      we think it can work some day.)
//   3. Javascript + DOM (done in content.js)
//   4. Marionette + DOM (doable)

//const originalWebSocket = window.wrappedJSObject.WebSocket;
var originalWebSocket = null;
var exportedWebsocket;

(function(){
  console.log("content.js loaded"); // are we really using our script?
  // can we have a visual clue that we are running our script?
  // this is not working, because the body style is overwritten by the website.
  //document.body.style.border = "5px solid red";

  // Save a reference to the original WebSocket constructor
  var nativeWebSocket = window.wrappedJSObject.WebSocket;

  // initialize connection to the quickBeam server.
  // This is the only connection that we will not intercept.
  var quickBeamSocket = null;
  const connectQuickBeam = () => {
    // it is important that we use the native WebSocket constructor here,
    // otherwise we would end up in an infinite loop.
    quickBeamSocket = new nativeWebSocket("ws://localhost:18981");

    // on disconnect, we will try to reconnect.
    quickBeamSocket.onclose = function (event) {
      console.error("quickBeamSocket closed", event);
      // we can just try to reconnect to the quickBeam server, firefox seems to
      // handle exponential backoff for us automatically. (does it really?)
      connectQuickBeam();
    };
  };

  // messageInterceptor is a function that will be called whenever a message
  // is received on any WebSocket connections, except our quickBeam socket.
  const messageInterceptor = (message) => {
    console.log("messageInterceptor", message);
    // forward the message to the quickBeam server.
    // we can deal with the message format server side.
    quickBeamSocket.send(message.data);
  }

  // We are wrapping the WebSocket constructor with our own class,
  // so we can intercept the calls and act on incoming messages.
  window.WebSocket = function (url, protocols) {
    // we are using the native WebSocket constructor for the heavy lifting,
    var that = protocols ? new nativeWebSocket(url, protocols) : new nativeWebSocket(url);
    // but we are intercepting the messages.
    that.addEventListener("message", messageInterceptor);

    // We don't care about the other events, since we assume that the connection
    // will be reestablished when there is a problem.
    that.addEventListener("open", console.log("socket open (intercepted)"));
    //that.addEventListener("close", console.log("socket close"));
    //that.addEventListener("error", console.log("socket error"));

    return that;
  };

  // finally we need to make sure that our implementation shares the same prototype,
  // so that there are absolutely no differences in behavior.
  window.WebSocket.prototype = nativeWebSocket.prototype;

  // don't know if this is necessary, but it can't hurt.
  originalWebSocket = window.WebSocket;

  // we push back the connection, so that the rest of the code has a chance to finish quicker.
  connectQuickBeam();

  console.log("content.js done"); // initialization done.
}());

// first i tried exporting the WebSocket constructor as a function, but that did not work.
//exportFunction(originalWebSocket, window, { defineAs: "WebSocket" });

// exporting the WebSocket constructor as an object works, but causes issues with the
// security model of Firefox. (xRayVision)
window.wrappedJSObject.WebSocket = cloneInto(originalWebSocket, window, {
  cloneFunctions: true,
});
