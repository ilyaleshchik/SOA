package api

type PaginationParams struct {
	Limit  int    `form:"limit"`
	PrevID string `form:"prev_id"`
}
