
const userQuerySelector = 'div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > span:nth-child(1)';
const messageQuerySelector = '[data-test="chatUserMessageText"]'


(function(){
    console.log("Loading extension ...");

    // Check for marker; only load the script once per page.
    if (document.head.querySelector('meta[name="quickbeam"]')) {
        console.warn('Script already loaded!');
        return; // abort early
    }
    const marker = document.createElement('meta');
    marker.setAttribute('name', 'quickbeam');
    document.head.appendChild(marker);

    // quickBeamSocket is the connection to the quickBeam server.
    var quickBeamSocket = null;

    const connectQuickBeam = () => {
      quickBeamSocket = new WebSocket("ws://localhost:18981");

      // on disconnect, we will try to reconnect.
      quickBeamSocket.onclose = function (event) {
        console.error("quickBeamSocket closed", event);
        // we can just try to reconnect to the quickBeam server, firefox seems to
        // handle exponential backoff for us automatically. (does it really?)
        connectQuickBeam();
      };
    };

    const handleChatMessage = (isPrivate, user, msg) => {
        scope = isPrivate ? "private" : "public";
        var message = new Object();
        message.type = "chat";
        message.scope = scope;
        message.user = user;
        message.message = msg;
        quickBeamSocket.send(JSON.stringify(message));
        //console.log(`(${scope}) ${user}: ${msg}`);
    };

    // we have 2 modes of operation:
    // - Public messages (default) are using the public toolbar
    // - Private messages are using the private tab

    const handleMessageNode = (node) => {
      const msgNode = node.querySelector(messageQuerySelector)
      // if the node was not found, it was not a chat message.
      if (!msgNode) {
        return
      }

      msg = msgNode.textContent;
      const user = addedNode.querySelector(userQuerySelector).textContent;
      handleChatMessage(false, user, msg);
    }

    const callback = (mutationsList, observer) => {
        for (const mutation of mutationsList) {
          if (mutation.type === "childList") {
              if (!mutation.addedNodes || mutation.addedNodes === []) {
                continue
              }
              for (const addedNode of mutation.addedNodes) {
                // check if element has attribute, in which case it is probably a chat message.
                // otherwise we can return early.
                if (addedNode.dataset.test === "msgListItem") handleMessageNode(addedNode);

              }
          } else if (mutation.type === "attributes") {
            // no content changes, just attribute changes.
            //console.log(`The ${mutation.attributeName} attribute was modified.`);
          }
        }
    };

    const observer = new MutationObserver(callback);
    observer.observe(document.body, { childList: true, subtree: true });
})();
