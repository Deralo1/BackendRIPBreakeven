package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Order struct { // вот наша новая структура
	ID               int    // поля структур, которые передаются в шаблон
	Title            string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
	ImageURL         string
	ShortDescription string
	Description      string
}

func (r *Repository) GetOrders() ([]Order, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	orders := []Order{
		{ID: 1,
			Title:            "Аренда склада",
			ImageURL:         "http://localhost:9000/lab1/arenda.jpg",
			ShortDescription: "Средняя цена за 1 кв.м в Москве: 1000 ₽",
			Description:      "Москва (в пределах МКАД): 800-1 200 ₽/м²  Московская область (до 20 км от МКАД): 500-800 ₽/м² Новая Москва: 600-900 ₽/м² Логистические парки МО: 450-700 ₽/м² Санкт-Петербург и ЛО СПб (промзоны): 500-800 ₽/м²  Ленинградская область (до 30 км): 350-550 ₽/м² Кудрово, Шушары: 450-650 ₽/м².",
		},
		{ID: 2, Title: "Оклад", ImageURL: "http://localhost:9000/lab1/oklad.jpg", ShortDescription: "Средний оклад по отрасли: 168 000 ₽", Description: "Заглушка: подробное описание услуги будет здесь."},
		{ID: 3, Title: "Сырье", ImageURL: "http://localhost:9000/lab1/siriy.jpg", ShortDescription: "Цены зависят от типа сырья", Description: "Заглушка: подробное описание услуги будет здесь."},
		{ID: 4, Title: "Сдельная Зарплата", ImageURL: "http://localhost:9000/lab1/sdelnay.jpg", ShortDescription: "Оплата за единицу продукции", Description: "Заглушка: подробное описание услуги будет здесь."},
		{ID: 5, Title: "Аренда торгового помещения", ImageURL: "http://localhost:9000/lab1/arendratorgpom.png", ShortDescription: "Средняя цена за 1 кв.м: 10000 ₽", Description: "Заглушка: подробное описание услуги будет здесь."},
		{ID: 6, Title: "Транспортировка", ImageURL: "http://localhost:9000/lab1/transportation.jpg", ShortDescription: "Средняя цена за 1 км: 20 ₽", Description: "Заглушка: подробное описание услуги будет здесь."},
	}
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
	if len(orders) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return orders, nil
}
func (r *Repository) GetOrder(id int) (Order, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Order{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}
func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}

	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}

	return result, nil
}
