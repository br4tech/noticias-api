package g1_scraper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/br4tech/noticias-api/internal/core/domain"
	"golang.org/x/net/html"
)

type scraper struct {
	client        *http.Client
	logger        *slog.Logger
	reURI         *regexp.Regexp
	reID          *regexp.Regexp
	reURL         *regexp.Regexp
	reTitle       *regexp.Regexp
	reDescription *regexp.Regexp
}

func NovoScraper(client *http.Client, logger *slog.Logger) *scraper {
	return &scraper{
		client:        client,
		logger:        logger,
		reURI:         regexp.MustCompile(`COMENTARIOS_URI:\s*"([^"]+)"`),
		reID:          regexp.MustCompile(`COMENTARIOS_IDEXTERNO:\s*"([^"]+)"`),
		reURL:         regexp.MustCompile(`CANONICAL_URL:\s*"([^"]+)"`),
		reTitle:       regexp.MustCompile(`TITLE:\s*"([^"]+)"`),
		reDescription: regexp.MustCompile(`DESCRIPTION:\s*"([^"]+)"`),
	}
}

func (s *scraper) BuscarSettingsDaNoticia(ctx context.Context, urlNoticia string) (*domain.SettingsNoticia, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlNoticia, nil)
	if err != nil {
		s.logger.Error("falha ao criar requisição http", "url", urlNoticia, "erro", err)
		return nil, fmt.Errorf("falha ao criar requisição: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("falha ao executar requisição http", "url", urlNoticia, "erro", err)
		return nil, fmt.Errorf("falha ao executar requisição: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("recebido status http não esperado", "url", urlNoticia, "status", resp.StatusCode)
		return nil, fmt.Errorf("status inesperado: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		s.logger.Error("falha ao analisar o html da página", "url", urlNoticia, "erro", err)
		return nil, fmt.Errorf("falha ao analisar html: %w", err)
	}

	scriptContent, found := s.findScriptSettings(doc)
	if !found {
		err := errors.New("a tag <script id='SETTINGS'> não foi encontrada na página")
		s.logger.Error(err.Error(), "url", urlNoticia)
		return nil, err
	}

	settings := &domain.SettingsNoticia{
		IDExterno:    s.extractString(s.reID, scriptContent),
		CanonicalURL: s.extractString(s.reURL, scriptContent),
		Title:        s.extractString(s.reTitle, scriptContent),
		Description:  s.extractString(s.reDescription, scriptContent),
	}

	if settings.Title == "" || settings.IDExterno == "" {
		return nil, errors.New("falha ao extrair dados do script, settings essenciais estão vazios")
	}

	return settings, nil
}

func (s *scraper) BuscarComentarios(ctx context.Context, settings *domain.SettingsNoticia, tipoOrdenacao string) (*domain.Noticia, error) {
	urlTemplate := "https://comentarios.globo.com/comentarios/{uri}/{idExterno}/{url}/{shorturl}/{titulo}/{pagina}.json"
	if tipoOrdenacao == "populares" {
		urlTemplate = "https://comentarios.globo.com/comentarios/{uri}/{idExterno}/{url}/{shorturl}/{titulo}/populares/{pagina}.json"
	}

	replacer := strings.NewReplacer(
		// "{uri}", url.QueryEscape(settings.URI),
		"{idExterno}", url.QueryEscape(settings.IDExterno),
		"{url}", url.QueryEscape(settings.CanonicalURL),
		"{shorturl}", "shorturl",
		"{titulo}", url.QueryEscape(settings.Title),
	)
	baseCommentURL := replacer.Replace(urlTemplate)

	var todosComentarios []domain.Comentario
	totalComentarios := 0
	pagina := 1
	const ItensPorPagina = 25

	for {
		requestURL := strings.Replace(baseCommentURL, "{pagina}", fmt.Sprintf("%d", pagina), 1)
		s.logger.Debug("Buscando comentários", "url", requestURL)

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
		if err != nil {
			return nil, fmt.Errorf("erro ao criar req de comentarios pagina %d: %w", pagina, err)
		}

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar comentarios pagina %d: %w", pagina, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status inesperado da API de comentarios pagina %d: %d", pagina, resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler corpo da resposta de comentarios pagina %d: %w", pagina, err)
		}

		jsonString := strings.TrimPrefix(string(bodyBytes), "__callback_listacomentarios(")
		jsonString = strings.TrimSuffix(jsonString, ")")

		var result struct {
			Itens []struct {
				Texto   string `json:"texto"`
				Replies []struct {
					Texto string `json:"texto"`
				} `json:"replies"`
			} `json:"itens"`
		}

		if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
			return nil, fmt.Errorf("erro ao decodificar JSON de comentarios pagina %d: %w", pagina, err)
		}

		for _, item := range result.Itens {
			respostas := make([]string, len(item.Replies))
			for i, reply := range item.Replies {
				respostas[i] = reply.Texto
			}
			todosComentarios = append(todosComentarios, domain.Comentario{
				Mensagem:  item.Texto,
				Respostas: respostas,
			})
			totalComentarios += 1 + len(respostas)
		}

		if len(result.Itens) < ItensPorPagina {
			break
		}

		pagina++
	}

	noticiaFinal := &domain.Noticia{
		Titulo: settings.Title,
		Total:  totalComentarios,
	}

	return noticiaFinal, nil
}

func (s *scraper) findScriptSettings(n *html.Node) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "script" {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == "SETTINGS" {
				if n.FirstChild != nil {
					return n.FirstChild.Data, true
				}
				return "", true
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if content, found := s.findScriptSettings(c); found {
			return content, true
		}
	}
	return "", false
}

func (s *scraper) extractString(re *regexp.Regexp, text string) string {
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
