package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/m21power/Ecom/service/cart"
	"github.com/m21power/Ecom/service/order"
	"github.com/m21power/Ecom/service/product"
	"github.com/m21power/Ecom/service/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	userStore := user.NewStore(s.db)
	// dependency injection
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)
	carthandler := cart.NewHandler(orderStore, productStore, userStore)
	carthandler.RegisterRoutes(subrouter)
	log.Println("listening on: ", s.addr)
	return http.ListenAndServe(s.addr, router)
}
