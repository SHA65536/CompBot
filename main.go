package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	bot := MakeCompBot(Config{
		Token:   os.Getenv("TOKEN"),
		Prefix:  os.Getenv("PREFIX"),
		Channel: os.Getenv("CHANNEL"),
	})
	bot.Start()

	// Gracefully shutting down
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	bot.Stop()
}
