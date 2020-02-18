package upload

// Response is a JSON response as a result of an upload
type Response struct {
	URL    string   `json:"url"`
	URLs   []string `json:"urls"`
	Errors []string `json:"errors"`
}

// PasteResponse is a JSON response as a result of a paste
type PasteResponse struct {
	URL    string   `json:"url"`
	Errors []string `json:"errors"`
}

// PastUpload is a JSON result as a result of a query of past uploads
type PastUpload struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	LocalName    string `json:"local_name"`
	Extension    string `json:"extension"`
	HasThumbnail bool   `json:"has_thumb"`
}

// PastUploads is a list of past uploads returned as a JSON response
type PastUploads struct {
	Files    []PastUpload `json:"files"`
	NumPages int          `json:"number_pages"`
	Errors   []string     `json:"errors"`
}

// DeleteResponse is a response sent as a result of a file/paste deletion
type DeleteResponse struct {
	Errors []string `json:"errors"`
}
