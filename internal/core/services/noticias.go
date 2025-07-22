package services

import (
	"context"
	"log/slog"
	"sync"

	"github.com/br4tech/noticias-api/internal/core/domain"
	"github.com/br4tech/noticias-api/internal/core/ports"
)

type servicoNoticia struct {
	feedRepo ports.FeedRepository
	scraper  ports.ScraperNoticias
	logger   *slog.Logger
}

func NovoServicoNoticias(feedRepo ports.FeedRepository, scraper ports.ScraperNoticias, logger *slog.Logger) ports.ServicoNoticias {
	return &servicoNoticia{feedRepo: feedRepo, scraper: scraper, logger: logger}
}

func (s *servicoNoticia) ObtereNoticias(ctx context.Context, categoria string, tipoOrdenacao string, limite int) ([]domain.SettingsNoticia, error) {

	urls, err := s.feedRepo.ObterURLsRecentes(categoria, limite)
	if err != nil {
		s.logger.Error("falha ao obter URLs recentes", "categoria", categoria, "erro", err)
		return nil, err
	}
	s.logger.Info("URLs obtidas com sucesso", "quantidade", len(urls))

	var wg sync.WaitGroup
	resultadosChan := make(chan domain.SettingsNoticia, len(urls))

	for _, urlNoticia := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			log := s.logger.With("url", url)
			log.Info("Iniciando scraping da notícia")

			settings, err := s.scraper.BuscarSettingsDaNoticia(ctx, url)
			if err != nil {
				log.Error("falha ao buscar settings da notícia", "erro", err)
				return
			}

			log.Info("Scraping da notícia concluído com sucesso", "titulo", settings.Title)
			resultadosChan <- *settings
		}(urlNoticia)
	}

	go func() {
		wg.Wait()
		close(resultadosChan)
	}()

	var settingsNoticiasFinais []domain.SettingsNoticia
	for noticia := range resultadosChan {
		settingsNoticiasFinais = append(settingsNoticiasFinais, noticia)
	}

	s.logger.Info("Processamento concorrente finalizado", "noticias_processadas", len(settingsNoticiasFinais))

	return settingsNoticiasFinais, nil
}

func (s *servicoNoticia) ObterNoticiaAleatoria(ctx context.Context, categoria string, tipoOrdenacao string) (*domain.Noticia, error) {

	urlNoticia, err := s.feedRepo.ObterURLNoticiaAleatoria(categoria)
	if err != nil {
		s.logger.Error("falha ao obter URL da notícia", "categoria", categoria, "erro", err)
		return nil, err
	}
	s.logger.Info("URL da notícia obtida com sucesso", "url", urlNoticia)

	settings, err := s.scraper.BuscarSettingsDaNoticia(ctx, urlNoticia)
	if err != nil {
		s.logger.Error("falha ao buscar settings da notícia", "url", urlNoticia, "erro", err)
		return nil, err
	}
	s.logger.Info("Settings da notícia extraídos", "titulo", settings.Title)

	noticia, err := s.scraper.BuscarComentarios(ctx, settings, tipoOrdenacao)
	if err != nil {
		s.logger.Error("falha ao buscar comentários", "titulo", settings.Title, "erro", err)
		return nil, err
	}
	s.logger.Info("Comentários buscados com sucesso", "total", noticia.Total)

	return noticia, nil
}
