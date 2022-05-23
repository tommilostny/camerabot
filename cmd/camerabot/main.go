package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/tommilostny/camerabot"
	_ "github.com/tommilostny/camerabot/handler"
)

const defaultCacheDir = "/tmp"

func main() {
	if v := os.Getenv("MAIN_CHAT_ID"); v != "" {
		var err error
		if camerabot.MainChatID, err = strconv.ParseInt(v, 10, 64); err != nil {
			log.Fatalf("Main chat ID is not a valid integer: %s", err)
		}
	} else {
		log.Fatal("MAIN_CHAT_ID env var not set")
	}

	if v := os.Getenv("CACHE_DIR"); v != "" {
		camerabot.CacheDir = v
	} else {
		log.Printf("Using default cache directory: %s", defaultCacheDir)
		camerabot.CacheDir = defaultCacheDir
	}

	if v := camerabot.ParseAllowedChatIDs(); v != nil {
		camerabot.AllowedChatIDs = v
	} else {
		log.Printf("No allowed chat IDs set. Allowing everyone.")
		camerabot.AllowedChatIDs = nil
	}

	s := gocron.NewScheduler(time.UTC)
	handler := camerabot.Handlers["/pic"]

	s.Every(1).Hour().Do(func() {
		for e := camerabot.AllowedChatIDs.Front(); e != nil; e = e.Next() {

			log.Println("Sending picture to chat:", e.Value)

			if err := handler.Handle(e.Value.(int64)); err != nil {
				log.Printf("Failed to handle chat ID %d: %s", e.Value.(int64), err)
			}
		}
	})

	go camerabot.ListenAndServe()

	s.StartAsync()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Println("Interrupt received. Graceful shutdown.")

	s.Stop()
}
