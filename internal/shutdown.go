package internal

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SetupGracefulShutdown(port string, engine *gin.Engine) {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	defer func() {
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown: ", err)
		}
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Application failed", err)
		}
	}()
	log.WithFields(log.Fields{"bind": port}).Info("Running application")

	WaitingForExitSignal()
	log.Info("Waiting for all jobs to stop")
}

func WaitingForExitSignal() {
	signalForExit := make(chan os.Signal, 1)
	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	stop := <-signalForExit
	log.Info("Stop signal Received", stop)
}
