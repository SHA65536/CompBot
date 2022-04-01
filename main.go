package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	cfg := Config{
		Token:      os.Getenv("TOKEN"),
		Channel:    os.Getenv("CHANNEL"),
		Prefix:     "!comp",
		CreateCD:   time.Minute * 5,
		ReactCD:    time.Second * 3,
		GameStatus: " CS:GO Since 1970",
	}
	if val, ok := os.LookupEnv("PREFIX"); ok {
		cfg.Prefix = val
	}
	if val, ok := os.LookupEnv("CREATE_CD"); ok {
		cfg.CreateCD, _ = time.ParseDuration(val)
	}
	if val, ok := os.LookupEnv("REACT_CD"); ok {
		cfg.ReactCD, _ = time.ParseDuration(val)
	}
	if val, ok := os.LookupEnv("GAME_STATUS"); ok {
		cfg.GameStatus = val
	}
	bot := MakeCompBot(cfg)
	bot.Start()

	// Gracefully shutting down
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	bot.Stop()
}
