package server

import (
	"context"
	"errors"
	"exampleproject/log"
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
	th := http.TimeoutHandler(handler, maxRequestDuration, "TimeoutHandler reached deadline")
	return &http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  30 * time.Second,
		Addr:         ":8080",
		Handler:      th,
	}
}

func Listen() (<-chan error, chan<- os.Signal) {
	handler := newHandler()
	server := newDefaultServer(handler)

	signalSigInt := make(chan os.Signal, 1)
	signal.Notify(signalSigInt, syscall.SIGINT)
	signalShutdown := make(chan error, 1)

	go func() {
		<-signalSigInt
		shutdownGracefully(server)
	}()

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.TEMPDDEBUG("gracefully stopped")
				signalShutdown <- nil
			} else {
				log.Log(log.Serve, log.Err(err))
				signalShutdown <- err
			}
		}
	}()

	return signalShutdown, signalSigInt
}

func shutdownGracefully(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), maxRequestDuration)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Log(log.ServerShutdown, log.Err(err))
		return
	}
}
