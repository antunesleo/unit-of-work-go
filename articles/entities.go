package articles

type Article struct {
	Id       int64    `json:"id"`
	Title    string   `json:"title"`
	Desc     string   `json:"desc"`
	Content  string   `json:"content"`
	Category Category `json:"category"`
}

type Category struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
