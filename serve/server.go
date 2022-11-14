package serve

import (
	"context"
	"errors"
	"exampleproject/log"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	maxRequestDuration = time.Second * 10
)

func newDefaultServer(handler http.Handler) *http.Server {
	return &http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  30 * time.Second,
		Addr:         ":8080",
		Handler:      http.TimeoutHandler(handler, maxRequestDuration, "TimeoutHandler reached deadline"),
	}
}

func StartServer() error {
	handler := newHandler()
	server := newDefaultServer(handler)

	onSigInt := make(chan os.Signal, 1)
	signal.Notify(onSigInt, syscall.SIGINT)
	onShutdownFinished := make(chan struct{}, 1)

	go func() {
		<-onSigInt
		shutdownGracefully(server, onShutdownFinished)
	}()

	err := server.ListenAndServe()
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			<-onShutdownFinished
			return nil
		}
		log.Log(log.Serve, log.Err(err))
		return err
	}

	return nil
}

func shutdownGracefully(server *http.Server, onShutdownFinished chan<- struct{}) {
	defer func() { onShutdownFinished <- struct{}{} }()
	ctx, cancel := context.WithTimeout(context.Background(), maxRequestDuration)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Log("shutdown error", log.Err(err))
		return
	}

	fmt.Println("gracefully stopped")
}
