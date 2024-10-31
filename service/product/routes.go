package product

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/m21power/Ecom/types"
	"github.com/m21power/Ecom/utils"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)

}
func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var p types.Product
	err := utils.ParseJSON(r, &p)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	err = h.store.CreateProduct(&p)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, p)
}
func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)

}
