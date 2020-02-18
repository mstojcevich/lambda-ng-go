function login(form) {
    let fData = new FormData(form);

    let request = new XMLHttpRequest();
    request.open("POST", "/api/user/login", true);

    request.onload = function () {
        // If the login was successful, go to the homepage
        if (request.status == 200) {
            window.location.href = "/"
        } else {
            // Login was unsuccessful, show the user the error
            let errorArea = document.getElementById("errorArea");
            errorArea.innerHTML = ""; // Clear out all errors
            try {
                let response = JSON.parse(request.responseText);
                console.log(response);
                if (response.errors.length > 0) {
                    response.errors.forEach(function (error) {
                        errorArea.innerHTML += "<div class=\"form-error\">" + error + "</div>";
                    });
                } else {
                    errorArea.innerHTML += "<div class=\"form-error\">Failed to login for unknown reason</div>";
                }
            } catch (e) {
                console.error(e);
                errorArea.innerHTML += "<div class=\"form-error\">Failed to login for unknown reason</div>";
            }
        }
    }

    // Send the login info
    request.send(fData);

    return false
}