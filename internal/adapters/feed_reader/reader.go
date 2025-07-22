package feed_reader

import (
	"context"
	"errors"

	"log/slog"
	"math/rand"
	"sync"

	"github.com/mmcdole/gofeed"
)

var feedsURL = map[string]string{
	"todas":            "http://pox.globo.com/rss/g1/",
	"brasil":           "http://pox.globo.com/rss/g1/brasil/",
	"carros":           "http://pox.globo.com/rss/g1/carros/",
	"ciencia_saude":    "http://pox.globo.com/rss/g1/ciencia-e-saude/",
	"economia":         "http://pox.globo.com/rss/g1/economia/",
	"educacao":         "http://pox.globo.com/rss/g1/educacao/",
	"loterias":         "http://pox.globo.com/rss/g1/loterias/",
	"mundo":            "http://pox.globo.com/rss/g1/mundo/",
	"planeta-bizarro":  "http://pox.globo.com/rss/g1/planeta-bizarro/",
	"politica":         "http://pox.globo.com/rss/g1/politica/",
	"pop-arte":         "http://pox.globo.com/rss/g1/pop-arte/",
	"tecnologia":       "http://pox.globo.com/rss/g1/tecnologia/",
	"turismo-e-viagem": "http://pox.globo.com/rss/g1/turismo-e-viagem/",
}

type feedRepository struct {
	cache  map[string][]string
	mu     sync.RWMutex
	logger *slog.Logger
}

func NovoFeedRepository(logger *slog.Logger) *feedRepository {
	return &feedRepository{
		cache:  make(map[string][]string),
		logger: logger,
	}
}

func (r *feedRepository) CarregarFeeds(ctx context.Context) error {
	r.logger.Info("Iniciando carregamento dos feeds RSS")
	fp := gofeed.NewParser()
	for cat, url := range feedsURL {
		feed, err := fp.ParseURLWithContext(url, ctx)
		if err != nil {
			r.logger.Warn("falha ao carregar feed", "categoria", cat, "erro", err)
			continue
		}

		r.mu.Lock()
		r.cache[cat] = []string{}
		for _, item := range feed.Items {
			r.cache[cat] = append(r.cache[cat], item.GUID)
		}
		r.mu.Unlock()
	}
	r.logger.Info("Feeds carregados e cacheados com sucesso")
	return nil
}

func (r *feedRepository) ObterURLNoticiaAleatoria(categoria string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	urls, ok := r.cache[categoria]
	if !ok || len(urls) == 0 {
		return "", errors.New("categoria não encontrada ou sem notícias: " + categoria)
	}

	return urls[rand.Intn(len(urls))], nil
}

func (r *feedRepository) ObterURLsRecentes(categoria string, limite int) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	urls, ok := r.cache[categoria]
	if !ok || len(urls) == 0 {
		return nil, errors.New("categoria não encontrada ou sem notícias: " + categoria)
	}

	if len(urls) < limite {
		return urls, nil
	}

	return urls[:limite], nil
}
