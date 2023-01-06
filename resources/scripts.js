let elementsUnderChange = {}

function ajaxCall(url) {
    const xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = () => {
        if (xhttp.readyState === 4 && xhttp.response !== undefined) {
            let modifiers = xhttp.response;
            if (modifiers == "")
                return;
            try {
                const obj = JSON.parse(modifiers);
                applyJSModifiers(obj.operations);
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
    if (!document.hidden)
        ajaxCall("/status")
}, 200);