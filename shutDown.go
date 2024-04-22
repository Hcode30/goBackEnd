package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func PrepareShutDown(server *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		<-c
		println("[WARNING]: Termination Signal Recieved!")
		res := saveAllData(data)
		if res {
			println("[INFO]: all data saved\n[WARNING]: Quiting...")
		} else {
			println("[ERROR]: Data is not saved!")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			saveAllData(data)
		}
	}()
}
