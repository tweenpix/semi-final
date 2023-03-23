package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"semifinal/lib"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	// Создание контекста для отслеживания сигналов ОС
	ctx, cancel := context.WithCancel(context.Background())

	// Создание HTTP-сервера
	server := &http.Server{Addr: ":9000"}
	service, err := lib.NewCityService("cities.csv")
	if err != nil {
		log.Fatal(err)
	}
	// Запуск сервера в горутине
	go func() {

		r := chi.NewRouter()
		r.Get("/{id}", service.GetCityByID)                     //получение информации о городе по его id;
		r.Post("/add/", service.AddCity)                        // добавление новой записи в список городов;
		r.Delete("/del/{id}", service.DeleteCityByID)           // удаление информации о городе по указанному id;
		r.Post("/update/{id}", service.UpdatePopulationByID)    // обновление информации о численности населения города по указанному id;
		r.Get("/cities-r/", service.GetCitiesByRegion)          // получение списка городов по указанному региону;
		r.Get("/cities-d/", service.GetCitiesByDistrict)        // получение списка городов по указанному округу;
		r.Get("/cities-p/", service.GetCitiesByPopulationRange) // получения списка городов по указанному диапазону численности населения;
		r.Get("/cities-f/", service.GetCitiesByFoundationRange) // получения списка городов по указанному диапазону года основания.

		http.ListenAndServe(":9000", r)

	}()

	// Ожидание сигнала от ОС
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	// Ожидание сигнала от ОС или отмены контекста
	select {
	case <-osSignal:
		fmt.Println("Received OS signal. Shutting down...")
		// сохранение данных в файл
		err := service.Save("cities.csv")
		if err != nil {
			log.Fatal(err)
		}
	case <-ctx.Done():
		fmt.Println("Context canceled. Shutting down...")
	}

	// Остановка HTTP-сервера
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
	}

	// Вызов функции отмены контекста
	cancel()
}
