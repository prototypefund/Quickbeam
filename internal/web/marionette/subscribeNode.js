console.log("subscribe");
(function (node, id) {
    console.log("subsribe to: ", node, " with id: ", id)
    const callback = (mutationList, observer) => {
        console.log(mutationList)
        console.log("subscription id: ", id)
        for (const mutation of mutationList) {
            if (mutation.type === "childList") {
                const msg = new Object();
                msg.type = "subscription";
                msg.id = id;
                msg.additions = mutation.addedNodes ? mutation.addedNodes.length : 0;
                msg.removals = mutation.removedNodes ? mutation.removedNodes.length : 0;
                console.log(msg)
                window.quickbeam.send(msg);
            }
        }
    };
    const observer = new MutationObserver(callback);
    observer.observe(node, { childList: true, subtree: false });
})(arguments[0], arguments[1]);
