package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kolya9390/RPC_JSON_GeoProvider/client_Proxy/router"
)

func main() {

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router.NewApiRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Создание канала для получения сигналов остановки
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Ожидание сигнала остановки
	<-stop

	log.Println("Shutting down server...")

	// Создание контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Остановка сервера с использованием graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped gracefully")

}