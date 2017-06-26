requiresLogin = true;

document.addEventListener("DOMContentLoaded", function() {
    getSessionInfo(function() {}, function() {
        // Redirect the user to login if they aren't signed in
        window.history.replaceState("backward", null, "/");
        window.location.href = "/login";
    });
});

function submitPaste() {
    let pastePlaintext = document.getElementById("paste-area").value;

    let encryptionKey = genEncKey(8);
    let encryptedPaste = encrypt(pastePlaintext, encryptionKey);
    
    let encryptedPasteObj = JSON.parse(encryptedPaste);
    encryptedPasteObj.is_code = false; // TODO have the code button do something
    encryptedPaste = JSON.stringify(encryptedPasteObj);

    putPaste(encryptedPaste, function(url) {
        window.location.href = "/" + url + "#" + encryptionKey;
    });
}

function putPaste(encryptedPaste, onFinish) {
    let request = new XMLHttpRequest();
    request.open("POST", "/api/paste", true);

    request.onload = function() {
        if(request.status == 200) { // success
            try {
                let response = JSON.parse(request.responseText);
                onFinish(response.url);
            } catch (e) {
                console.error(e);
                alert("An unexpected error occurred");
            }
        } else {
            try {
                let response = JSON.parse(request.responseText);
                response.errors.forEach(function(error) {
                    alert(error);
                });
            } catch (e) {
                console.error(e);
                alert("An unexpected error occurred");
            }
        }
    }

    let data = new FormData();
    data.append("paste", encryptedPaste);

    request.send(data);
}

/**
 * Generates an encryption key
 */
function genEncKey(length) {
    let charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    
    let key = "";
    while(key.length < length) {
        key += charset.charAt(Math.floor(Math.random() * charset.length));
    }

    return key;
}

function encrypt(text, key) {
    return sjcl.encrypt(key, text);
}