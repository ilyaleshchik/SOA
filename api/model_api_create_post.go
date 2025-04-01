package api

type CreatePost struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Private     bool     `json:"is_private"`
	Tags        []string `json:"tags"`
}

type PostUpdate struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Private     *bool     `json:"is_private"`
	Tags        *[]string `json:"tags"`
}

type Post struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Private     bool     `json:"is_private"`
	Tags        []string `json:"tags"`
	OwnerID     int64    `json:"owner_id"`
	DateCreate  int64    `json:"date_create"`
}
