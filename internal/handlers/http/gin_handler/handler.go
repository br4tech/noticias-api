package gin_handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/br4tech/noticias-api/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	servico ports.ServicoNoticias
	logger  *slog.Logger
}

func NovoHandler(servico ports.ServicoNoticias, logger *slog.Logger) *Handler {
	return &Handler{servico: servico, logger: logger}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	// router.GET("/favicon.ico", func(c *gin.Context) {
	// 	c.Status(http.StatusNoContent)
	// })

	router.GET("/", h.obterNoticiasHandler)
	router.GET("/:categoria", h.obterNoticiasHandler)
	router.GET("/populares", h.obterNoticiasHandler)
	router.GET("/:categoria/populares", h.obterNoticiasHandler)
}

func (h *Handler) obterNoticiasHandler(c *gin.Context) {
	categoria := strings.ToLower(c.Param("categoria"))
	if categoria == "" {
		categoria = "todas"
	}

	tipoOrdenacao := "recentes"
	if strings.Contains(c.FullPath(), "/populares") { // Verificação mais robusta
		tipoOrdenacao = "populares"
	}

	const limiteDeNoticias = 10

	log := h.logger.With("categoria", categoria, "ordenação", tipoOrdenacao, "limite", limiteDeNoticias)
	log.Info("recebida nova requisição para múltiplas notícias")

	noticias, err := h.servico.ObtereNoticias(c.Request.Context(), categoria, tipoOrdenacao, limiteDeNoticias)
	if err != nil {
		log.Error("erro no serviço", "erro", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao processar a requisição"})
		return
	}

	c.JSON(http.StatusOK, noticias)
}
