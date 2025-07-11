package bff

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luizdequeiroz/rpg-backend/internal/app/models"
	"github.com/luizdequeiroz/rpg-backend/internal/app/services"
)

// SheetTemplateHandler gerencia endpoints de templates de ficha
type SheetTemplateHandler struct {
	service *services.SheetTemplateService
}

// NewSheetTemplateHandler cria uma nova instância do handler
func NewSheetTemplateHandler(service *services.SheetTemplateService) *SheetTemplateHandler {
	return &SheetTemplateHandler{
		service: service,
	}
}

// ListTemplates godoc
// @Summary Lista templates de ficha
// @Description Retorna uma lista paginada de todos os templates de ficha ativos
// @Tags templates
// @Accept json
// @Produce json
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Itens por página" default(20)
// @Success 200 {object} models.SheetTemplateListResponse "Lista de templates"
// @Failure 400 {object} models.ErrorResponse "Parâmetros inválidos"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/templates [get]
func (h *SheetTemplateHandler) ListTemplates(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	response, err := h.service.GetAll(page, limit)
	if err != nil {
		if err == services.ErrInvalidPagination {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Parâmetros de paginação inválidos",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Erro interno do servidor",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetTemplate godoc
// @Summary Busca template por ID
// @Description Retorna os detalhes de um template específico
// @Tags templates
// @Accept json
// @Produce json
// @Param id path int true "ID do template"
// @Success 200 {object} models.SheetTemplateResponse "Detalhes do template"
// @Failure 400 {object} models.ErrorResponse "ID inválido"
// @Failure 404 {object} models.ErrorResponse "Template não encontrado"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/templates/{id} [get]
func (h *SheetTemplateHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "ID inválido",
		})
		return
	}

	template, err := h.service.GetByID(id)
	if err != nil {
		if err == services.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Template não encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Erro interno do servidor",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// CreateTemplate godoc
// @Summary Cria um novo template
// @Description Cria um novo template de ficha de personagem
// @Tags templates
// @Accept json
// @Produce json
// @Param template body models.CreateSheetTemplateRequest true "Dados do template"
// @Success 201 {object} models.SheetTemplateResponse "Template criado com sucesso"
// @Failure 400 {object} models.SheetTemplateErrorResponse "Dados inválidos"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/templates [post]
func (h *SheetTemplateHandler) CreateTemplate(c *gin.Context) {
	var req models.CreateSheetTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "JSON inválido",
			Message: err.Error(),
		})
		return
	}

	// Validar dados
	if validationErrors := h.service.ValidateCreateRequest(req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, models.SheetTemplateErrorResponse{
			BaseErrorResponse: models.BaseErrorResponse{
				Error:     "Dados inválidos",
				Timestamp: time.Now(),
			},
			Fields: validationErrors,
		})
		return
	}

	template, err := h.service.Create(req)
	if err != nil {
		if err == services.ErrInvalidDefinition {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Definição JSON inválida",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Erro interno do servidor",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// UpdateTemplate godoc
// @Summary Atualiza um template
// @Description Atualiza os dados de um template existente
// @Tags templates
// @Accept json
// @Produce json
// @Param id path int true "ID do template"
// @Param template body models.UpdateSheetTemplateRequest true "Dados atualizados do template"
// @Success 200 {object} models.SheetTemplateResponse "Template atualizado com sucesso"
// @Failure 400 {object} models.SheetTemplateErrorResponse "Dados inválidos"
// @Failure 404 {object} models.ErrorResponse "Template não encontrado"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/templates/{id} [put]
func (h *SheetTemplateHandler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "ID inválido",
		})
		return
	}

	var req models.UpdateSheetTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "JSON inválido",
			Message: err.Error(),
		})
		return
	}

	// Validar dados
	if validationErrors := h.service.ValidateUpdateRequest(req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, models.SheetTemplateErrorResponse{
			BaseErrorResponse: models.BaseErrorResponse{
				Error:     "Dados inválidos",
				Timestamp: time.Now(),
			},
			Fields: validationErrors,
		})
		return
	}

	template, err := h.service.Update(id, req)
	if err != nil {
		if err == services.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Template não encontrado",
			})
			return
		}

		if err == services.ErrInvalidDefinition {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Definição JSON inválida",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Erro interno do servidor",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, template)
}

// DeleteTemplate godoc
// @Summary Remove um template
// @Description Remove um template de ficha (soft delete)
// @Tags templates
// @Accept json
// @Produce json
// @Param id path int true "ID do template"
// @Success 204 "Template removido com sucesso"
// @Failure 400 {object} models.ErrorResponse "ID inválido"
// @Failure 404 {object} models.ErrorResponse "Template não encontrado"
// @Failure 500 {object} models.ErrorResponse "Erro interno do servidor"
// @Router /api/v1/templates/{id} [delete]
func (h *SheetTemplateHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "ID inválido",
		})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if err == services.ErrTemplateNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Template não encontrado",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Erro interno do servidor",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// SetupTemplateRoutes configura as rotas de templates
func (h *SheetTemplateHandler) SetupTemplateRoutes(router *gin.RouterGroup) {
	templates := router.Group("/templates")
	{
		templates.GET("", h.ListTemplates)
		templates.GET("/:id", h.GetTemplate)
		templates.POST("", h.CreateTemplate)
		templates.PUT("/:id", h.UpdateTemplate)
		templates.DELETE("/:id", h.DeleteTemplate)
	}
}
