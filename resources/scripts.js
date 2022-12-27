function ajaxCall(url) {
    const xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = () => {
        if (xhttp.readyState === 4 && xhttp.response !== undefined) {
            let modifiers = xhttp.response;
            applyJSModifiers(modifiers);
        }
    }
    xhttp.open('GET', url, true);
    xhttp.send();
}

function applyJSModifiers(modifiers) {
    // šel by použít rovnou eval, ale to je security zlo
    let modifiersArr = modifiers.split(";")
    for (var i = 0; i < modifiersArr.length; i++) {
        let modifier = modifiersArr[i];
        let vars = modifier.split(",")
        if (vars.length > 1) {
            if (vars[0] == "addClass") {
                document.getElementById(vars[1]).classList.add(vars[2]);
            } else if (vars[0] == "removeClass") {
                document.getElementById(vars[1]).classList.remove(vars[2]);
            } else if (vars[0] == "showError") {
                let infoDiv = document.getElementById("info-div");
                infoDiv.innerText = vars[1]
            } else if (vars[0] == "songInfo") {
                let infoDiv = document.getElementById("current-song-div");
                infoDiv.innerText = vars[1]
            }
        }
    }
}

setInterval(() => {
    ajaxCall("/status")
}, 500);