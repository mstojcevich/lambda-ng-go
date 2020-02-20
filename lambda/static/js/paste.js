requiresLogin = true;

var isCode = false;

document.addEventListener("DOMContentLoaded", function() {
    getSessionInfo(function() {}, function() {
        // Redirect the user to login if they aren't signed in
        window.history.replaceState("backward", null, "/");
        window.location.href = "/login";
    });
    registerEventHandlers();
});

function registerEventHandlers() {
    document.getElementById("codeToggleBtn").addEventListener("click", toggleCode);
    document.getElementById("submitBtn").addEventListener("click", submitPaste);
}

function submitPaste(e) {
    e.preventDefault();

    let pastePlaintext = document.getElementById("paste-area").value;

    let encryptionKey = genEncKey(8);
    let encryptedPaste = encrypt(pastePlaintext, encryptionKey);
    
    let encryptedPasteObj = JSON.parse(encryptedPaste);
    encryptedPasteObj.is_code = isCode;
    encryptedPaste = JSON.stringify(encryptedPasteObj);

    putPaste(encryptedPaste, isCode, function(url) {
        window.location.href = "/" + url + "#" + encryptionKey;
    });
}

function putPaste(encryptedPaste, isCode, onFinish) {
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
    data.append("is_code", isCode);

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

function toggleCode(e) {
    e.preventDefault();
    if(!isCode) {
        isCode = true;
        document.getElementById("codeLabel").innerText = "CODE: YES";
    } else {
        isCode = false;
        document.getElementById("codeLabel").innerText = "CODE: NO";
    }
}