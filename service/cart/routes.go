package cart

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/m21power/Ecom/service/auth"
	"github.com/m21power/Ecom/types"
	"github.com/m21power/Ecom/utils"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{
		store:        store,
		productStore: productStore,
		userStore:    userStore,
	}
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)

}
func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	var cart types.CartCheckoutPayload
	userID := auth.GetIdFromContext(r.Context())
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload : %v", errors))
		return
	}
	// get products

	prIDs, err := GetCartItemsIDs(cart.Item)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ps, err := h.productStore.GetProductByIDs(prIDs)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	orderID, totalPrice, err := h.createOrder(ps, cart.Item, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})

}
