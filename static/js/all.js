var sessionInfo = null
var attemptedAuth = false

var requiresLogin = false // Whether the current page requires login

/**
 * Gets info about the current Lambda session
 * @param {*} onSuccess Function called when session data is successfully obtained. It is passed a variable w/ session info.
 * @param {*} onFail Function called when session data could not be obtained
 */
function getSessionInfo(onSuccess, onFail) {
    if(!attemptedAuth) {
        let request = new XMLHttpRequest();
        request.open("GET", "/api/session", true);

        request.onerror = onFail;

        request.onload = function () {
            // If the request was successful, call onSuccess w/ the info obtained
            if (request.status == 200) {
                try {
                    sessionInfo = JSON.parse(request.responseText);
                    onSuccess(sessionInfo);
                    attemptedAuth = true;
                } catch (ex) { // Didn't recieve valid json. This shouldn't happen.
                    console.error(ex);
                    onFail(response);
                    attemptedAuth = true;
                }
            } else {
                try {
                    // The user probably wasn't signed in
                    response = JSON.parse(request.responseText);
                    onFail(response)
                    attemptedAuth = true;
                } catch (ex) {
                    // Other error. No idea what happened to get us here.
                    console.error(ex);
                    onFail()
                    attemptedAuth = true;
                }
            }
        }

        // Make the request to the server
        request.send();
    } else {
        if(sessionInfo != null) {
            onSuccess(sessionInfo);
        } else {
            onFail();
        }
    }
}