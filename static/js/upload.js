requiresLogin = true;

const uploadURL = "/api/upload";
const uploadDomain = "/";

document.addEventListener("DOMContentLoaded", function() {
    getSessionInfo(function() {}, function() {
        // Redirect the user to login if they aren't
        window.history.replaceState("backward", null, "/");
        window.location.href = "/login";
    });

    // On file hover
    document.body.addEventListener("dragover", function(e) {
        e.stopPropagation();
        e.preventDefault();
        e.dataTransfer.dropEffect = "copy";

        // TODO mdn mentions a "drag image". Look into that
    });

    // On file drop
    document.body.addEventListener("drop", function(e) {
        e.stopPropagation();
        e.preventDefault(); // Prevent browser redirect

        var files = e.dataTransfer.files;
        console.log(files);
        for(let i = 0; i < files.length; i++) {
            console.log(files[i]);

            // TODO properly handle multiple uploads in one request
            checkAndUpload(files[i]);
        }
    });

    let selectInput = document.getElementById("chooseFile");
    selectInput.addEventListener("change", function() {
        // TODO properly handle multiple uploads
        checkAndUpload(selectInput.files[0]);
    });
});

/**
 * Checks to see that the user is logged in
 * and the file meets all requirements
 * then proceeds to upload file
 * @param {*} file File to upload
 */
function checkAndUpload(file) {
    let errorList = document.getElementById("errorList");

    getSessionInfo(function() { // User is logged in
        if(file.size > maxFilesize) {
            addError(errorList, "Max filesize is " + (maxFilesize / (1024 * 1024)) + " MB");
            return
        }

        // TODO check extension
        upload(file, onUploadFinish);
    }, function() { // User isn't logged in
        addError(errorList, "You must be logged in to upload files");
    });
}

function upload(file, onFinish) {
    let xhr = new XMLHttpRequest();
    let fd = new FormData();
    
    xhr.open("PUT", uploadURL, true);
    fd.append("file", file);

    // Add the status indicator showing upload progress
    createStatusIndicator(xhr, file);

    // Handle upload finish
    xhr.onreadystatechange = function() {
        if(xhr.readyState == 4) { // Request completed
            if(xhr.status == 200) { // Request successful
                let response = JSON.parse(xhr.responseText);
                response.file = file; // Attach the file so that it can be referred to in the onFinish function
                onFinish(response);
            } else {
                try {
                    // Decode the response then show errors to the user
                    let response = JSON.parse(xhr.responseText);
                    response.errors.forEach(function(error) {
                        addError(errorList, error);
                    });
                } catch(e) {
                    console.error(xhr.responseText);
                    console.error(e);
                    addError(errorList, "An unknown error occurred");
                }
            }
        }
    }

    xhr.send(fd);
}

function createStatusIndicator(xhr, file) {
    let ongoingSection = document.getElementById("ongoing-uploads");
    let uploadEl = document.createElement("li");

    // If the file is an image then add an image preview to the status indicator
    if(isImage(file)) {
        let image = document.createElement("img");
        image.src = URL.createObjectURL(file);
        uploadEl.appendChild(image);
    }

    let uploadContent = document.createElement("div");

    let title = document.createElement("span");
    title.innerText = file.name;

    let progress = document.createElement("progress");
    progress.value = 0;

    uploadContent.appendChild(title);
    uploadContent.appendChild(progress);
    uploadEl.appendChild(uploadContent);
    ongoingSection.appendChild(uploadEl);

    // Keep the progress bar updated
    xhr.upload.addEventListener("progress", function(e) {
        let percent = e.loaded / e.total;
        progress.value = percent;

        // Remove upload on finish
        if(percent >= 1.0) {
            ongoingSection.removeChild(uploadEl);
        }
    });
}

function addError(errorList, error) {
    let li = document.createElement("li");
    li.innerText = error;
    errorList.appendChild(li);
}

function onUploadFinish(response) {
    let url = response.url;
    
    let finishedSection = document.getElementById("finished-uploads");
    let uploadEl = document.createElement("li");

    // Clicking a finished upload should bring the user to the uploaded file
    uploadEl.onclick = function() {
        window.location = uploadDomain + url;
    }

    // If the file is an image then add an image preview to the status indicator
    if(isImage(response.file)) {
        let image = document.createElement("img");
        image.src = URL.createObjectURL(response.file);
        uploadEl.appendChild(image);
    }

    let contentArea = document.createElement("div");

    let uploadLink = document.createElement("a");
    uploadLink.href = uploadDomain + url;

    uploadLink.innerText = response.file.name;

    contentArea.appendChild(uploadLink);
    uploadEl.appendChild(contentArea);
    finishedSection.appendChild(uploadEl);
}

function isImage(file) {
    return file.type.lastIndexOf("image/", 0) == 0  // beginsWith('image/')
}