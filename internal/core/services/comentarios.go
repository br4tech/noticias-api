package services

import (
	"context"
	"log/slog"

	"github.com/br4tech/noticias-api/internal/core/domain"
	"github.com/br4tech/noticias-api/internal/core/ports"
)

type servicoComentarios struct {
	feedRepo ports.FeedRepository
	scraper  ports.ScraperNoticias
	logger   *slog.Logger
}

func NovoServicoComentarios(feedRepo ports.FeedRepository, scraper ports.ScraperNoticias, logger *slog.Logger) ports.ServicoComentarios {
	return &servicoComentarios{feedRepo: feedRepo, scraper: scraper, logger: logger}
}

func (s *servicoComentarios) ObterComentariosDeNoticia(ctx context.Context, categoria string, tipoOrdenacao string) (*domain.Noticia, error) {

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
