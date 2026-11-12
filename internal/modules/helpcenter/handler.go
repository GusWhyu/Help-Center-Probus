package helpcenter

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafialfauzan/layrd/pkg/database"
	"github.com/rafialfauzan/layrd/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Failed to retrieve categories", err.Error())
		return
	}
	response.OK(c, "Categories retrieved successfully", categories)
}

func (h *Handler) GetCategoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid category ID", nil)
		return
	}

	category, err := h.service.GetCategoryByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NotFound(c, "Category not found")
			return
		}
		response.InternalError(c, "Failed to retrieve category", err.Error())
		return
	}
	response.OK(c, "Category retrieved successfully", category)
}

func (h *Handler) GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	kategoriID, _ := strconv.ParseUint(c.Query("kategori_id"), 10, 32)
	menuID, _ := strconv.ParseUint(c.Query("menu_id"), 10, 32)
	search := c.Query("q")

	articles, total, err := h.service.GetArticles(c.Request.Context(), page, limit, uint(kategoriID), uint(menuID), search)
	if err != nil {
		response.InternalError(c, "Failed to retrieve articles", err.Error())
		return
	}

	response.Paginated(c, "Articles retrieved successfully", articles, page, limit, total)
}

func (h *Handler) GetArticleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid article ID", nil)
		return
	}

	article, err := h.service.GetArticleByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			response.NotFound(c, "Article not found")
			return
		}
		response.InternalError(c, "Failed to retrieve article", err.Error())
		return
	}
	response.OK(c, "Article retrieved successfully", article)
}

func (h *Handler) GetManualBooks(c *gin.Context) {
	menuID, _ := strconv.ParseUint(c.Query("menu_id"), 10, 32)
	search := c.Query("q")

	manuals, err := h.service.GetManualBooks(c.Request.Context(), uint(menuID), search)
	if err != nil {
		response.InternalError(c, "Failed to retrieve manual books", err.Error())
		return
	}
	response.OK(c, "Manual books retrieved successfully", manuals)
}

func (h *Handler) GetManualBookByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid manual book ID", nil)
		return
	}

	manual, err := h.service.GetManualBookByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			response.NotFound(c, "Manual book not found")
			return
		}
		response.InternalError(c, "Failed to retrieve manual book", err.Error())
		return
	}
	response.OK(c, "Manual book retrieved successfully", manual)
}

func (h *Handler) Search(c *gin.Context) {
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	results, err := h.service.Search(c.Request.Context(), query, limit)
	if err != nil {
		response.InternalError(c, "Failed to search help center", err.Error())
		return
	}
	response.OK(c, "Search results retrieved successfully", results)
}
