document.addEventListener("DOMContentLoaded", function() {
    registerEventHandlers();
});

function registerEventHandlers() {
    document.getElementById("loginForm").addEventListener("submit", login);
}

function login(e) {
    e.preventDefault();
    let fData = new FormData(e.target);

    let request = new XMLHttpRequest();
    request.open("POST", "/api/user/login", true);

    request.onload = function () {
        // If the login was successful, go to the homepage
        if (request.status == 200) {
            window.location.href = "/";
        } else {
            // Login was unsuccessful, show the user the error
            let errorArea = document.getElementById("errorArea");
            errorArea.textContent = ""; // Clear out all errors
            try {
                let response = JSON.parse(request.responseText);
                if (response.errors.length > 0) {
                    response.errors.forEach(function (error) {
                        let errorDiv = document.createElement("div");
                        errorDiv.className = "form-error";
                        errorDiv.textContent = error;
                        errorArea.appendChild(errorDiv);
                    });
                } else {
                    let errorDiv = document.createElement("div");
                    errorDiv.className = "form-error";
                    errorDiv.textContent = "Failed to login for unknown reason";
                    errorArea.appendChild(errorDiv);
                }
            } catch (e) {
                console.error(e);
                let errorDiv = document.createElement("div");
                errorDiv.className = "form-error";
                errorDiv.textContent = "Failed to login for unknown reason";
                errorArea.appendChild(errorDiv);
            }
        }
    }

    // Send the login info
    request.send(fData);
}