package domain

type Comentario struct {
	Mensagem  string   `json:"mensagem"`
	Respostas []string `json:"respostas"`
}

type Noticia struct {
	Titulo string `json:"titulo"`
	Total  int    `json:"total"`
}

type SettingsNoticia struct {
	IDExterno    string `json:"id"`
	CanonicalURL string `json:"url"`
	Title        string `json:"title"`
	SubTitle     string `json:"sub_title"`
	Description  string `json:"description"`
}
