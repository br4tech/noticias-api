package domain

// Comentario representa um comentário com suas respostas.
type Comentario struct {
	Mensagem  string   `json:"mensagem"`
	Respostas []string `json:"respostas"`
}

// Noticia representa o resultado final que será retornado pela API.
type Noticia struct {
	Titulo      string       `json:"titulo"`
	Total       int          `json:"total"`
	Comentarios []Comentario `json:"comentarios"`
}

// SettingsNoticia armazena os metadados extraídos do HTML de uma notícia.
type SettingsNoticia struct {
	URI          string
	IDExterno    string
	CanonicalURL string
	Title        string
}
