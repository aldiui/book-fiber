package dto

type MediaData struct {
	Id        string `json:"id"`
	Path      string `json:"path"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type CreateMediaRequest struct {
	Path string `json:"path"`
}
