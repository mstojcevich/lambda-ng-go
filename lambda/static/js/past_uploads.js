requiresLogin = true;

const perPage = 21;

var page = 1;

var selectedImages = [];

document.addEventListener("DOMContentLoaded", function() {
    registerEventHandlers();

    getUploads();

    window.addEventListener("hashchange", function() {
        getUploads();
    });
});

function registerEventHandlers() {
    document.getElementById("goBack").addEventListener("click", prevPage);
    document.getElementById("goNext").addEventListener("click", nextPage);
    document.getElementById("deleteSelectedBtn").addEventListener("click", deleteSelected);
}

function getUploads() {
    // Set the selected page based on the URL
    let specifiedPage = window.location.hash.substr(1);
    if(specifiedPage.length > 0) {
        page = specifiedPage;
    }

    // Clear out the past uploads to make room for new ones
    let uploads = document.getElementById("uploads");
    uploads.innerHTML = "";

    // Make a request to get past uploads
    let xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() {
        if(xmlHttp.readyState == 4) { // Finished
            if(xmlHttp.status == 200) { // success
                let response = JSON.parse(xmlHttp.responseText);

                // Update the page count
                let numPages = response.number_pages;
                let pageNum = document.getElementById("pageNum");
                pageNum.innerText = "Page " + page + " of " + numPages;

                // Hide the go back button if there's nowhere to go back
                let goBackButton = document.getElementById("goBack");
                if(page > 1) {
                    goBackButton.className = "navButton";
                } else {
                    goBackButton.className = "navButton hidden";
                }

                // Hide the next button if there's nowhere to go next
                let goFwdButton = document.getElementById("goNext");
                if(page < numPages) {
                    goFwdButton.className = "navButton";
                } else {
                    goFwdButton.className = "navButton hidden";
                }

                response.files.forEach(function(upload) {
                    let li = document.createElement("li");
                    li.title = upload.local_name;

                    let a = document.createElement("a");
                    a.href = "/" + upload.name + "." + upload.extension;

                    let img = document.createElement("img");
                    img.name = upload.name;

                    // Scale the image properly
                    img.onload = function(e) {
                        let width = e.target.clientWidth;
                        let height = e.target.clientHeight;

                        if(height > width) {
                            e.target.className = "tall";
                        } else {
                            e.target.className = "wide";
                        }
                    }

                    img.onerror = function(e) { // If the image failed to load
                        if(!img.erroredBefore) {
                            img.erroredBefore = true; // Prevent endless loop of loading the replacement image if the replacement image also fails
                            img.src = "/generic/by-ext/" + upload.extension;
                        }
                    }

                    // TODO image types that don't have thumbnails should still show
                    if(upload.has_thumb) {
                        // TODO hardcoded thumbnail url
                        img.src = "/thumb_128x128_" + upload.name + ".jpg";
                    } else {
                        img.src = "/generic/by-ext/" + upload.extension;
                    }

                    // Set alttext to the image's name
                    img.alt = upload.local_name;

                    a.appendChild(img);
                    li.appendChild(a);
                    uploads.appendChild(li);

                    // Handle selecting on right click
                    img.myLi = li;
                    img.oncontextmenu = function(e) {
                        toggleSelection(e.target, e.target.myLi);
                        return false;
                    }
                });
            } else if(xmlHttp.status == 401) { // Not logged in
                window.history.replaceState("backward", null, "/");
                window.location.href = "/login";
            } else { // failure
                response = JSON.parse(xmlHttp.responseText);
                console.error(response);
            }
        }
    }

    let reqUrl = "/api/user/uploads?page=" + page + "&n=" + perPage;
    xmlHttp.open("GET", reqUrl, true);
    xmlHttp.send();
}

function prevPage(e) {
    e.preventDefault();
    page = parseInt(page) - 1;
    window.location.hash = page;
}

function nextPage(e) {
    e.preventDefault();
    page = parseInt(page) + 1;
    window.location.hash = page;
}

function toggleSelection(img, li) {
    if(li.className.length == 0) {
        li.className = "selected";
        selectedImages.push(img);
    } else {
        li.className = "";
        // Remove img from selectedImages
        selectedImages.splice(selectedImages.indexOf(img), 1);
    }

    if(selectedImages.length == 0) {
        document.getElementById("selection-management").className = "selection-manage hidden";
    } else {
        document.getElementById("selection-management").className = "selection-manage";
        document.getElementById("numSelectedLabel").innerText = selectedImages.length + " items selected";
    }
}

function deleteSelected(e) {
    e.preventDefault();

    let toFinishCount = selectedImages.length;

    selectedImages.forEach(function(img) {
        deleteImage(img.name, function() {
            toFinishCount--;
            if(toFinishCount == 0) {
                location.reload(); // Refresh the page TODO just update the list
            }
        });
    });
}

function deleteImage(name, callback) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function(e) {
        if (e.target.readyState === 4) callback();
    };
    xmlHttp.open("DELETE", "/file/" + name, true);
    xmlHttp.send();
}
