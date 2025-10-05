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

type Service struct { // вот наша новая структура
	ID               int    // поля структур, которые передаются в шаблон
	Title            string // ОБЯЗАТЕЛЬНО должны быть написаны с заглавной буквы (то есть публичными)
	ImageURL         string
	Price            int
	ShortDescription string
	Description      string
}

func (r *Repository) GetCostsService() ([]Service, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	services :=
		[]Service{
			{
				ID:               1,
				Title:            "Аренда склада",
				ImageURL:         "http://localhost:9000/lab1/arenda.jpg",
				Price:            1000,
				ShortDescription: "Средняя стоимость аренды склада в Москве",
				Description:      "Аренда складских помещений востребована в Москве и области. Стоимость зависит от расположения: в пределах МКАД — 800–1200 ₽/м², в Подмосковье — 500–800 ₽/м². В логистических парках цена ниже, но склады находятся дальше от центра. При выборе учитывают транспортную доступность, наличие охраны и инфраструктуры.",
			},
			{
				ID:               2,
				Title:            "Оклад",
				ImageURL:         "http://localhost:9000/lab1/oklad.jpg",
				Price:            168000,
				ShortDescription: "Средний оклад по отрасли",
				Description:      "Оклад — фиксированная часть заработной платы, которая выплачивается сотруднику независимо от объёма выполненной работы. В среднем по отрасли он составляет около 168 000 ₽. Конкретная сумма зависит от должности, региона и уровня квалификации.",
			},
			{
				ID:               3,
				Title:            "Сырьё",
				ImageURL:         "http://localhost:9000/lab1/siriy.jpg",
				Price:            1000,
				ShortDescription: "Стоимость сырья для производства 1кг продукции",
				Description:      "Сырьё — это материалы, используемые для производства продукции. Цена зависит от вида: металл, древесина, пластик, текстиль и т.д. На стоимость влияют мировые котировки, логистика и объём закупки.",
			},
			{
				ID:               4,
				Title:            "Сдельная зарплата",
				ImageURL:         "http://localhost:9000/lab1/sdelnay.jpg",
				Price:            10,
				ShortDescription: "Оплата за единицу продукции",
				Description:      "Сдельная форма оплаты труда предполагает, что работник получает деньги за каждую произведённую единицу продукции или выполненную операцию. Такой подход стимулирует производительность, но требует строгого контроля качества.",
			},
			{
				ID:               5,
				Title:            "Аренда торгового помещения",
				ImageURL:         "http://localhost:9000/lab1/arendratorgpom.png",
				Price:            10000,
				ShortDescription: "Аренда торговых площадей в Москве",
				Description:      "Аренда торговых помещений востребована в крупных городах. Цена зависит от расположения: в центре Москвы — от 10 000 ₽/м², в спальных районах — дешевле. Важны проходимость, наличие парковки и транспортная доступность.",
			},
			{
				ID:               6,
				Title:            "Транспортировка",
				ImageURL:         "http://localhost:9000/lab1/transportation.jpg",
				Price:            20,
				ShortDescription: "Стоимость за киллометр транспортировки",
				Description:      "Транспортировка включает доставку товаров автомобильным, железнодорожным или морским транспортом. Средняя цена — 20 ₽ за километр. На стоимость влияют расстояние, вес груза и выбранный вид транспорта.",
			},
		}

	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
	if len(services) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return services, nil
}
func (r *Repository) GetService(id int) (Service, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	services, err := r.GetCostsService()
	if err != nil {
		return Service{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, order := range services {
		if order.ID == id {
			return order, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Service{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}
func (r *Repository) GetServicesByTitle(title string) ([]Service, error) {
	orders, err := r.GetCostsService()
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
