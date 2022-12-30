function ajaxCall(url) {
    const xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = () => {
        if (xhttp.readyState === 4 && xhttp.response !== undefined) {
            let modifiers = xhttp.response;
            const obj = JSON.parse(modifiers);
            applyJSModifiers(obj.operations);
        }
    }
    xhttp.open('GET', url, true);
    xhttp.send();
}

function applyJSModifiers(operations) {
    // šel by použít rovnou eval, ale to je security zlo
    for (var i = 0; i < operations.length; i++) {
        let operation = operations[i];
        let name = operation.name;
        let params = operation.params;
        switch (name) {
            case "addClass":
                document.getElementById(params[0]).classList.add(params[1]);
                break;
            case "removeClass":
                document.getElementById(params[0]).classList.remove(params[1]);
                break;
            case "showError":
                document.getElementById("info-div").innerText = params[0];
                break;
            case "songInfo":
                let infoDiv = document.getElementById("current-song-div").innerText = params[0]
                break;
        }
    }
}

setInterval(() => {
    ajaxCall("/status")
}, 500);