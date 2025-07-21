package ports

import (
	"context"

	"github.com/br4tech/noticias-api/internal/core/domain"
)

type FeedRepository interface {
	CarregarFeeds(ctx context.Context) error
	ObterURLNoticiaAleatoria(categoria string) (string, error)
}

type ScraperNoticias interface {
	BuscarSettingsDaNoticia(ctx context.Context, urlNoticia string) (*domain.SettingsNoticia, error)
	BuscarComentarios(ctx context.Context, settings *domain.SettingsNoticia, tipoOrdenacao string) (*domain.Noticia, error)
}

type ServicoComentarios interface {
	ObterComentariosDeNoticia(ctx context.Context, categoria string, tipoOrdenacao string) (*domain.Noticia, error)
}
