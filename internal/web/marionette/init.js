(function (port){
    console.log("Initializing quickbeam javascript ...");

    // Check for marker; only load the script once per page.
    if (window.quickbeam && window.quickbeam.loaded) {
        console.warn('Script already loaded!');
        return; // abort early
    }
    window.quickbeam = new Object();
    window.quickbeam.loaded = true;

    // window.quickbeam.socket holds the backchannel to quickbeam
    // it is connected to a websocket server on the go side
    window.quickbeam.socket = null;

    // window.quickbeam.send sends msg to quickbeam via the websocket channel
    window.quickbeam.send = (msg) => {
        console.log("window.quickbeam.send(msg): ", msg)
        if (window.quickbeam.socket) {
            window.quickbeam.socket.send(JSON.stringify(msg));
        } else {
            console.error("cannot send msg because socket is null");
        }
    };

    // connectQuickBeam connects to the websocket server
    const connectQuickBeam = () => {
        window.quickbeam.socket = new WebSocket(`ws://localhost:${port}`);

        // on disconnect, we will try to reconnect.
        window.quickbeam.socket.onclose = function (event) {
            console.error("window.quickbeam.socket closed", event);
            // we can just try to reconnect to the quickBeam server, firefox seems to
            // handle exponential backoff for us automatically. (does it really?)
            connectQuickBeam();
        };

    };
    connectQuickBeam();

})(arguments[0]);
