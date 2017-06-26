document.addEventListener("DOMContentLoaded", function() {
    if(window.location.hash.length > 0) {
        let pasteName = window.location.pathname.split("#")[0].split("/")[1];
        getPaste(pasteName);
    } else {
        document.getElementById("paste-area").innerText = "Unable to find encryption key for paste";
    }
});

function getPaste(name) {
    let request = new XMLHttpRequest();
    request.open("GET", "/api/paste?name=" + name, true);

    request.onload = function() {
        if(request.status == 200) {
            try {
                decryptPaste(request.response);
            } catch(ex) {
                console.error(ex);
                alert("An unexpected error occurred");
            }
        }
    }

    request.send();
}

function decryptPaste(encPaste) {
    let key = window.location.hash.substr(1);
    let json = JSON.parse(encPaste);

    // Get rid of the is_code marker so sjcl can handle the json
    let sjclJson = JSON.parse(encPaste);
    delete(sjclJson.is_code);
    let sjclPlain = JSON.stringify(sjclJson);

    document.getElementById("paste-area").innerText = sjcl.decrypt(key, sjclPlain);

    hljs.highlightBlock(document.getElementById("paste-area"));
}