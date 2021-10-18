package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"SewingWorkshop/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	var products []*models.Product
	if r.URL.Query().Has("type") {
		pType := r.URL.Query().Get("type")
		s, err := app.products.LatestWithType(pType)
		if err != nil {
			app.serverError(w, err)
			return
		}
		products = s
	} else {
		s, err := app.products.Latest()
		if err != nil {
			app.serverError(w, err)
			return
		}
		products = s
	}

	types, err := app.products.GetTypes()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Используем помощника render() для отображения шаблона.
	app.render(w, r, "home.page.tmpl", &templateData{
		Products:     products,
		ProductTypes: types,
	})
}

func (app *application) showCustomerOrders(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.products.GetOrdersForCustomer(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// Используем помощника render() для отображения шаблона.
	app.render(w, r, "show.page.tmpl", &templateData{
		Products: s,
	})
}

func (app *application) showOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.products.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Используем помощника render() для отображения шаблона.
	app.render(w, r, "show.page.tmpl", &templateData{
		Product: s,
	})
}

func (app *application) createOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Создаем несколько переменных, содержащих тестовые данные. Мы удалим их позже.
	tp := "Куртка"
	cost := 3000
	size := "S"
	material := "Хлопок"
	master := 1
	customer := 1

	// Передаем данные в метод ProductModel.Insert(), получая обратно
	// ID только что созданной записи в базу данных.
	id, err := app.products.Insert(tp, size, material, cost, master, customer)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Перенаправляем пользователя на соответствующую страницу заметки.
	http.Redirect(w, r, fmt.Sprintf("/order?id=%d", id), http.StatusSeeOther)
}
