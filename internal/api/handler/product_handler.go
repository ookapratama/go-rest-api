package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ookapratama/go-rest-api/internal/domain"
	"github.com/ookapratama/go-rest-api/internal/util"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		slog.Error("Failed to decode request body", "error", err)
		util.WriteError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := util.ValidateStruct(&p); err != nil {
		slog.Warn("Validation failed", "error", err)
		util.WriteError(w, http.StatusBadRequest, "Validation error", err)
		return
	}

	if err := h.service.Create(r.Context(), &p); err != nil {
		slog.Error("Failed to create product", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "Failed to create product", err)
		return
	}

	util.WriteSuccess(w, http.StatusCreated, "Product created successfully", p)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	p, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		slog.Error("Product not found", "id", id, "error", err)
		util.WriteError(w, http.StatusNotFound, "Product not found", err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, "Product retrieved", p)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Demonstrating concurrency with GetEnrichedAll
	products, err := h.service.GetEnrichedAll(r.Context())
	if err != nil {
		slog.Error("Failed to get products", "error", err)
		util.WriteError(w, http.StatusInternalServerError, "Failed to get products", err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, "Products retrieved", products)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}
	p.ID = id

	if err := util.ValidateStruct(&p); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Validation error", err)
		return
	}

	if err := h.service.Update(r.Context(), &p); err != nil {
		slog.Error("Failed to update product", "id", id, "error", err)
		util.WriteError(w, http.StatusInternalServerError, "Failed to update product", err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, "Product updated successfully", p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		slog.Error("Failed to delete product", "id", id, "error", err)
		util.WriteError(w, http.StatusInternalServerError, "Failed to delete product", err)
		return
	}

	util.WriteSuccess(w, http.StatusOK, "Product deleted successfully", nil)
}
