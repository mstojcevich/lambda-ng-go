document.addEventListener("DOMContentLoaded", function() {
    registerEventHandlers();
});

function registerEventHandlers() {
    document.getElementById("registerForm").addEventListener("submit", register);
}

function register(e) {
    e.preventDefault();
    let fData = new FormData(e.target);

    let request = new XMLHttpRequest();
    request.open("POST", "/api/user/new", true);

    request.onload = function() {
        if(request.status == 200) { // Successful registration
            window.location.href = "/";
        } else { // Unsuccessful
            try {
                let response = JSON.parse(request.responseText);
                if(response.errors.length > 0) {
                    let errorArea = document.getElementById("errorArea");
                    errorArea.textContent = ""; // Clear the error list

                    response.errors.forEach(function(error) {
                        let errorDiv = document.createElement("div");
                        errorDiv.className = "form-error";
                        errorDiv.textContent = error;
                        errorArea.appendChild(errorDiv);
                    });
                }
            } catch(ex) {
                console.error(ex);
                alert("Register failed for unknown reason");
            }

            console.error(request.responseText);
        }
    }

    request.send(fData);
}