package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/miaozhang/webservice/logging"
	"github.com/miaozhang/webservice/models"
	"github.com/miaozhang/webservice/routers"
	"github.com/miaozhang/webservice/settings"
)

func init() {
	settings.Setup()
	models.Setup()
	logging.Setup()
}

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%d", settings.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    settings.ServerSetting.ReadTimeout,
		WriteTimeout:   settings.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	//s.ListenAndServe()
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
