package api

type Reminder struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
