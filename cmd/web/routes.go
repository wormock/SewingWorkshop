package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/order", app.showOrder)
	mux.HandleFunc("/order/create", app.createOrder)
	mux.HandleFunc("/customer", app.showCustomerOrders)
	mux.HandleFunc("/master/delete", app.deleteMaster)
	mux.HandleFunc("/master/add", app.addMaster)
	mux.HandleFunc("/master", app.showMasters)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
