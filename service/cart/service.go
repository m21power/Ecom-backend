package cart

import (
	"fmt"
	"log"

	"github.com/m21power/Ecom/types"
)

func GetCartItemsIDs(items []types.CartItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the products %d ", item.ProductID)

		}
		productIds[i] = item.ProductID
	}
	return productIds, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}
	//check if all products are in stock
	if err := checkIfCartInStock(items, productMap); err != nil {
		return 0, 0, err
	}
	//calculate total price
	totalPrice := calculateTotalPrice(items, productMap)
	//reduce quantity of the product in our db
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		if err := h.productStore.UpdateProduct(product); err != nil {
			return 0, 0, err
		}
	}
	// create order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})
	if err != nil {
		log.Println("error creating order", err)
		return 0, 0, nil
	}
	// create order item
	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}
	return orderID, totalPrice, nil

}
func checkIfCartInStock(cartItems []types.CartItem, productMap map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart item is empty")
	}
	for _, item := range cartItems {
		product, ok := productMap[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", product.Name)
		}
	}
	return nil
}

func calculateTotalPrice(items []types.CartItem, productMap map[int]types.Product) float64 {
	var total float64
	for _, item := range items {
		product := productMap[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}
	return total
}
