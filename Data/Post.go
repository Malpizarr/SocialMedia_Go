package data

type Post struct {
	ID       string   `json:"id"`
	Content  string   `json:"content"`
	Likes    int      `json:"likes"`
	Comments []string `json:"comments"`
	ImageURL string   `json:"imageURL"`
}
