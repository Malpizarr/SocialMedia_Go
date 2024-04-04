package data

type Post struct {
	Content  string
	Likes    int
	Comments []string
	ImageURL string // Campo para almacenar la URL de la imagen cargada
}
