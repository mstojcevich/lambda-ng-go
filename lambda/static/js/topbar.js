var userDropdownToggled = false;

// Wait for DOM to be loaded before doing stuff that will update the DOM
document.addEventListener("DOMContentLoaded", function () {
    // Update the authentication part of the topbar
    getSessionInfo(function (sessionInfo) {
        // Successful login, switch the login button to be a user details button

        let topbarAccount = document.getElementById("topbar-account");
        topbarAccount.innerText = sessionInfo.username;

        // Have the username button show the user dropdown
        topbarAccount.href = "#"; // Remove link to login
        topbarAccount.onclick = toggleUserDropdown; // Show the user dropdown when clicked
    }, function () {
        document.getElementById("topbar-account").innerText = "Not Signed In";
    });
});

/**
 * Toggles the dropdown w/ additional user info
 */
function toggleUserDropdown() {
    let userDropdown = document.getElementById("user-dropdown");

    if (!userDropdownToggled) { // Show the user dropdown
        userDropdown.className = "user-dropdown";
        userDropdown.style.right = (document.documentElement.clientWidth - document.getElementById("topbar-account").getBoundingClientRect().right) + "px";
        userDropdownToggled = true;
    } else { // hide the user dropdown
        userDropdown.className = "user-dropdown hidden";
        userDropdownToggled = false;
    }
}

// From https://stackoverflow.com/questions/5968196/check-cookie-if-cookie-exists
function getCookie(name) {
    var dc = document.cookie;
    var prefix = name + "=";
    var begin = dc.indexOf("; " + prefix);
    if (begin == -1) {
        begin = dc.indexOf(prefix);
        if (begin != 0) return null;
    }
    else {
        begin += 2;
        var end = document.cookie.indexOf(";", begin);
        if (end == -1) {
            end = dc.length;
        }
    }
    // because unescape has been deprecated, replaced with decodeURI
    //return unescape(dc.substring(begin + prefix.length, end));
    return decodeURI(dc.substring(begin + prefix.length, end));
}

// Signs the user out
function signOut() {
    let request = new XMLHttpRequest();
    request.open("DELETE", "/api/session");

    let accountArea = document.getElementById("topbar-account");

    accountArea.innerText = "Signing out...";

    request.onload = function () {
        if (request.status == 200) { // Success
            // Update the topbar to reflect being signed out
            accountArea.innerText = "Not Signed In";
            accountArea.href = "/login";
            accountArea.onclick = null;

            // Hide the user dropdown
            document.getElementById("user-dropdown").className = "user-dropdown hidden";
            userDropdownToggled = false;

            // If the current page requires login then go to login page
            if(requiresLogin) {
                window.history.replaceState("backward", null, "/");
                window.location.href = "/login";
            }
        } else {
            console.error(request.responseText);
            accountArea.innerText = "Failed to sign out";
        }
    }

    request.send();
}
