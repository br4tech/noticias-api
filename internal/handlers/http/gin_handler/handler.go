package gin_handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/br4tech/noticias-api/internal/core/ports"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	servico ports.ServicoComentarios
	logger  *slog.Logger
}

func NovoHandler(servico ports.ServicoComentarios, logger *slog.Logger) *Handler {
	return &Handler{servico: servico, logger: logger}
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.GET("/", h.obterComentariosHandler)
	router.GET("/:categoria", h.obterComentariosHandler)
	router.GET("/populares", h.obterComentariosHandler)
	router.GET("/:categoria/populares", h.obterComentariosHandler)
}

func (h *Handler) obterComentariosHandler(c *gin.Context) {
	categoria := strings.ToLower(c.Param("categoria"))
	if categoria == "" {
		categoria = "todas"
	}

	tipoOrdenacao := "recentes"
	if strings.HasSuffix(c.FullPath(), "/populares") {
		tipoOrdenacao = "populares"
	}

	log := h.logger.With("categoria", categoria, "ordenação", tipoOrdenacao)
	log.Info("recebida nova requisição")

	noticia, err := h.servico.ObterComentariosDeNoticia(c.Request.Context(), categoria, tipoOrdenacao)
	if err != nil {
		log.Error("erro no serviço", "erro", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao processar a requisição"})
		return
	}

	c.JSON(http.StatusOK, noticia)
}
