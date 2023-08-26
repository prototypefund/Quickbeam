

(function(){
    const userQuerySelector = 'div:nth-child(1) > div:nth-child(2) > div:nth-child(1) > div:nth-child(1) > span:nth-child(1)';
    const messageQuerySelector = '[data-test="chatUserMessageText"]'

    const settingScript = document.querySelector("html > body > script:not([src])");
    const settingsJSON = settingScript.text.match(/__meteor_runtime_config__ = JSON.parse\(decodeURIComponent\("(.*)"\)\)/)[1];
    const settingsMeteorRuntimeConfig = JSON.parse(decodeURIComponent(settingsJSON));
    window.quickbeam.bbbAppSettings = settingsMeteorRuntimeConfig.PUBLIC_SETTINGS.app;

    // quickBeamSocket is the connection to the quickBeam server.
    var quickBeamSocket = null;

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
                if (addedNode.dataset && addedNode.dataset.test === "msgListItem") handleMessageNode(addedNode);

              }
          } else if (mutation.type === "attributes") {
            // no content changes, just attribute changes.
            //console.log(`The ${mutation.attributeName} attribute was modified.`);
          }
        }
    };

    const observer = new MutationObserver(callback);
    observer.observe(document.documentElement, { childList: true, subtree: true });
})();
