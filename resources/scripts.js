let elementsUnderChange = {};
let lastPlaylistHash = 0;

function ajaxCall(url, callback) {
    const xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = () => {
        if (xhttp.readyState === 4 && xhttp.response !== undefined) {
            if (xhttp.response == "" || typeof callback === "undefined")
                return;
            try {
                callback(JSON.parse(xhttp.response));
            } catch (error) {
                console.error(error);
            }
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
                document.getElementById("current-song-div").innerText = params[0];
                document.title = "GrassControl: " + params[0];
                break;
            case "volume":
                document.getElementById("volume-span").innerText = Math.floor(params[0] * 100 / 256) + "%";
                let volumeSliderId = "volume-slider";
                if (!elementsUnderChange[volumeSliderId])
                    document.getElementById(volumeSliderId).value = Number(params[0]);
                break;
            case "progress":
                let progressSliderId = "progress-slider";
                let time = Number(params[0])
                let length = Number(params[1])
                if (!elementsUnderChange[progressSliderId]) {
                    document.getElementById(progressSliderId).value = time;
                    document.getElementById(progressSliderId).max = length;
                    document.getElementById("progress-length-span").innerText = formatTime(length);
                    document.getElementById("progress-time-span").innerText = formatTime(time);
                }
                break;
            case "playlistSelect":
                let id = params[0]
                const className = "table-tr-selected";
                const collection = document.getElementsByClassName(className);
                let noChange = false;
                for (element of collection) {
                    // pokud je záznam již označen, nic neřeš
                    if (element.classList.contains(className) && element.id == id) {
                        noChange = true;
                        break;
                    }
                    element.classList.remove(className);
                }
                if (noChange)
                    break;
                let target = document.getElementById(id);
                if (target) {
                    target.scrollIntoView();
                    target.classList.add(className);
                }
                break;
        }
    }
}

function formatTime(time) {
    let m = Math.floor(time / 60);
    let s = time % 60;
    return (m > 9 ? m : "0" + m) + ":" + (s > 9 ? s : "0" + s)
}

function progressControlScroll(event, callback) {
    let slider = event.target;
    let newVal = Number(slider.value) + Math.sign(-event.deltaY) * 5;
    slider.value = newVal;
    callback(newVal);
}

function volumeControlScroll(event, callback) {
    let slider = event.target;
    let newVal = Math.min(320, Number(slider.value) + Math.sign(-event.deltaY) * 5);
    slider.value = newVal;
    callback(newVal);
}

setInterval(() => {
    if (!document.hidden) {
        ajaxCall("/status", json => {
            applyJSModifiers(json.operations);
        });
        ajaxCall("/playlist", json => {
            hash = Number(json["hash"])
            if (hash != lastPlaylistHash) {
                lastPlaylistHash = hash;
                html = json["html"]
                document.getElementById("playlist-table-div").innerHTML = html;
            }
        });
    }
}, 200);