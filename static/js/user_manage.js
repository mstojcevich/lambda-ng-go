requiresLogin = true;

document.addEventListener("DOMContentLoaded", function () {
    // Update the authentication part of the topbar
    getSessionInfo(function (sessionInfo) {
        // Successful login, display the API key

        let topbarAccount = document.getElementById("topbar-account");
        console.log(sessionInfo);
        document.getElementById("apiKey").innerText = sessionInfo.api_key;

    }, function () {
        document.getElementById("apiKey").innerText = "Not Signed In";
    });
});
