package handler

import "net/http"

func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/orders", CreateOrderHandler)
}
