function ajaxCall(url, callback) {
    const xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = () => {
        if (xhttp.readyState === 4 && callback !== undefined)
            callback()
    }
    xhttp.open('GET', url, true);
    xhttp.send();
}