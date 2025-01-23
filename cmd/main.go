package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/config"
	mylog "url-shortener/internal/lib/log"
	"url-shortener/internal/service"
	"url-shortener/internal/service/handlers"
	"url-shortener/internal/storage/postgres"

	"github.com/julienschmidt/httprouter"
)

func main() {
	cnf := config.LoadConfig();

	log := mylog.LoadLogger(cnf.ENV);

	log.Debug("Init logger", slog.String("env", cnf.ENV));

	storage, err := postgres.New(cnf.Postgres);

	if err != nil {
		log.Error("Error connect POSTGRES", slog.String("error", err.Error()));
		os.Exit(1);
	}

	log.Debug("Init storage");

	app := service.New(*cnf, storage, log);

	log.Debug("Init app");

	router := httprouter.New();

	router.POST("/Savealias", handlers.SaveUrl(app));
	router.GET("/Getalias/:alias", handlers.GetAlias(app));


	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port),
		Handler: handlers.RecoverPanic(app, router),
		IdleTimeout: cnf.Server.IdleTimeout,
		ReadTimeout: cnf.Server.Timeout,
		WriteTimeout: cnf.Server.Timeout,
	};

	// для сигналов остановки сервера
	done := make(chan os.Signal, 1);
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM);

	go func(){
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", 
					slog.String("error:", err.Error()));
		}
	} ()

	log.Info("server started");

	//TODO: Закрыть сервер при сигналах
	<-done //примем значение, если будет сигнал
	log.Info("stopping server");

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second);
	defer cancel();

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server");
		return;
	}

}