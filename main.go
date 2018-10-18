package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/youyo/shiftscheduler/routers"
)

const (
	GracefulShutdownTimeout time.Duration = 30 * time.Second
)

var (
	ListenAddr string = os.Getenv("ADDR") + ":" + os.Getenv("PORT")
)

func main() {
	router := routers.Setup()

	srv := &http.Server{
		Addr:    ListenAddr,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
