package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	apiURL         = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1" // URL для получения данных о криптовалютах
	updateInterval = 10 * time.Minute                                                                                           // Интервал обновления данных
)

// Структура для хранения данных о криптовалюте
type Currency struct {
	ID     string  `json:"id"`
	Symbol string  `json:"symbol"`
	Name   string  `json:"name"`
	Price  float64 `json:"current_price"`
}

// Глобальные переменные для хранения данных и времени последнего обновления
var (
	currencies  []Currency
	lastUpdated time.Time
	mux         sync.Mutex // Мьютекс для синхронизации доступа к данным
)

// Функция для получения данных о криптовалютах из API
func fetchCurrencies() ([]Currency, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var curr []Currency
	if err := json.NewDecoder(resp.Body).Decode(&curr); err != nil {
		return nil, err
	}

	return curr, nil
}

// Проверяем, нужно ли обновить данные
func updateCurrenciesPeriodically() {
	updateData()
	for {
		time.Sleep(updateInterval)
		updateData()
	}
}

// Функция для обновления данных
func updateData() {
	mux.Lock()
	defer mux.Unlock()
	var err error
	currencies, err = fetchCurrencies()
	if err != nil {
		log.Printf("Ошибка при получении данных о криптовалютах: %v", err)
	} else {
		lastUpdated = time.Now()
	}
}

// Обработчик для получения списка всех криптовалют
func getCurrencies(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()

	json.NewEncoder(w).Encode(currencies)
}

// Обработчик для получения данных о конкретной криптовалюте
func getCurrencyByID(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	defer mux.Unlock()

	// Извлекаем ID из URL
	id := r.URL.Path[len("/currency/"):]

	for _, currency := range currencies {
		if currency.ID == id {
			json.NewEncoder(w).Encode(currency)
			return
		}
	}

	http.Error(w, "Криптовалюта не найдена", http.StatusNotFound)
}

func main() {

	go updateCurrenciesPeriodically()

	// Регистрируем обработчики маршрутов
	http.HandleFunc("/currencies", getCurrencies)
	http.HandleFunc("/currency/", getCurrencyByID)

	log.Println("Сервер запущен на порту :8080")
	http.ListenAndServe(":8080", nil)
}
