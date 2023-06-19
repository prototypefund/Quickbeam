var clicks = 0;

function addAndRemove() {
    var target = document.getElementById("changingList");

    clicks++;
    switch (clicks) {
    case 1:
        target.innerHTML += "<li>New Item 1</li>";
        break;
    case 2:
        target.innerHTML += "<li>New Item 2</li>";
        break;
    case 3:
        var victim = target.children[0];
        target.removeChild(victim);
        break;
    case 4:
        target.innerHTML += "<li>New Item 3</li>";
        break;
    case 5:
        var victim = target.children[2];
        target.removeChild(victim);
        break;
    }
}
