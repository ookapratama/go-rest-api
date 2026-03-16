package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ookapratama/go-rest-api/latihan_mandiri/internal/domain"
	"github.com/ookapratama/go-rest-api/latihan_mandiri/internal/util"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	
	// 1. Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		util.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 2. Call Service Create
	err := h.service.Create(r.Context(), &p)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 3. Respon Success
	util.WriteJSON(w, http.StatusCreated, "Product created", p)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Panggil GetEnrichedAll untuk mencoba Concurrency
	products, err := h.service.GetEnrichedAll(r.Context())
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.WriteJSON(w, http.StatusOK, "Products retrieved", products)
}
