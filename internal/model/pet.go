package model

type Pet struct {
	ID        int      `json:"id" example:"1"`
	Category  Category `json:"category"`
	Name      string   `json:"name" example:"Rex"`
	PhotoUrls []string `json:"photoUrls" example:"[\"https://example.com/photo.jpg\"]"`
	Tags      []Tag    `json:"tags"`
	Status    string   `json:"status" example:"available"`
}

type Category struct {
	ID   int    `json:"id" example:"2"`
	Name string `json:"name" example:"Dog"`
}

type Tag struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"cute"`
}
