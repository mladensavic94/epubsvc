package main

import (
	"os"
	"os/signal"
	"syscall"

	epubsvc "github.com/mladensavic94/epubsvc/internal"
)

func main() {
	epubsvc.InitLogger("")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go epubsvc.Server(":8080", epubsvc.NewStorage())
	epubsvc.Logger.Info("epubsvc started on :8080")
	<-c
	epubsvc.Logger.Info("epubsvc exited!")
}
