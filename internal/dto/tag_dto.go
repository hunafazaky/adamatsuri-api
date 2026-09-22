package dto

// TagResponse is the public shape of a tag — used standalone (GET /tags)
// and nested inside EventResponse.
type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
