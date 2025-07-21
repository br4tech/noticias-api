package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/br4tech/noticias-api/internal/adapters/feed_reader"
	"github.com/br4tech/noticias-api/internal/adapters/g1_scraper"
	"github.com/br4tech/noticias-api/internal/core/services"
	"github.com/br4tech/noticias-api/internal/handlers/http/gin_handler"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Configurar Observabilidade (Logger Estruturado)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 2. Inicializar os Adaptadores
	httpClient := &http.Client{Timeout: 10 * time.Second}
	feedRepo := feed_reader.NovoFeedRepository(logger)
	scraper := g1_scraper.NovoScraper(httpClient, logger)

	// Carregar o cache de feeds na inicialização
	if err := feedRepo.CarregarFeeds(context.Background()); err != nil {
		logger.Error("falha fatal ao carregar feeds", "erro", err)
		os.Exit(1)
	}

	// 3. Inicializar o Serviço do Core (injetando os adaptadores)
	servicoComentarios := services.NovoServicoComentarios(feedRepo, scraper, logger)

	// 4. Inicializar o Adaptador de Entrada (HTTP Handler)
	router := gin.Default()
	handler := gin_handler.NovoHandler(servicoComentarios, logger)
	handler.SetupRoutes(router)

	// 5. Iniciar o Servidor
	logger.Info("servidor iniciando na porta :8080")
	if err := router.Run(":8080"); err != nil {
		logger.Error("falha ao iniciar o servidor http", "erro", err)
		os.Exit(1)
	}
}
