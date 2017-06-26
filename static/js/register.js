function register() {
    let registerForm = document.getElementById(registerForm);
    let fData = new FormData(registerForm);

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
                    errorArea.innerHTML = ""; // Clear the error list

                    response.errors.forEach(function(error) {
                        errorArea.innerHTML += "<div class=\"form-error\">" + error + "</div>";
                    });
                }
            } catch(ex) {
                console.error(ex);
                alert("Register failed for unknown reason");
            }

            console.error(request.responseText);
        }
    }
}