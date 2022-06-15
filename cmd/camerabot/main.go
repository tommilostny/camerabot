package main

import (
	"container/list"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
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

	if v := parseIntListFromEnvironment("ALLOWED_CHAT_IDS"); v != nil {
		camerabot.AllowedChatIDs = v
		log.Printf("Allowed chat IDs:")
		for e := camerabot.AllowedChatIDs.Front(); e != nil; e = e.Next() {
			log.Printf("  %d", e.Value.(int64))
		}
	} else {
		log.Printf("No allowed chat IDs set. Allowing everyone.")
		camerabot.AllowedChatIDs = nil
	}

	if v := parseIntListFromEnvironment("MINUTES"); v != nil {
		camerabot.Minutes = v
		log.Printf("Minutes set:")
		for e := camerabot.Minutes.Front(); e != nil; e = e.Next() {
			log.Printf("  %d", e.Value.(int64))
		}
	} else {
		log.Printf("No minutes to auto-send set. Auto-send disabled.")
		camerabot.Minutes = nil
	}

	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Minute().Do(sendPhotoToAllAllowedChatIDs)

	go camerabot.ListenAndServe()

	s.StartAsync()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Println("Interrupt received. Graceful shutdown.")

	s.Stop()
}

func parseIntListFromEnvironment(envVarName string) *list.List {
	splitted := strings.Split(os.Getenv(envVarName), ";")
	if len(splitted) == 0 {
		return nil
	}

	result := list.New()
	for i := 0; i < len(splitted); i++ {
		value, err := strconv.ParseInt(splitted[i], 0, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		result.PushBack(value)
	}

	if result.Len() == 0 {
		return nil
	}
	return result
}

func shouldSendAPhoto(minute int64) bool {
	if camerabot.Minutes != nil {
		for e := camerabot.Minutes.Front(); e != nil; e = e.Next() {
			if e.Value.(int64) == minute {
				return true
			}
		}
	}
	return false
}

func sendPhotoToAllAllowedChatIDs() {
	if !shouldSendAPhoto(int64((61 + time.Now().Minute()) % 60)) {
		return
	}
	time.Sleep(time.Second * time.Duration(60-time.Now().Second()))

	handler := camerabot.Handlers["/pic"]
	for e := camerabot.AllowedChatIDs.Front(); e != nil; e = e.Next() {

		log.Println("Sending picture to chat:", e.Value)

		if err := handler.Handle(e.Value.(int64)); err != nil {
			log.Printf("Failed to handle chat ID %d: %s", e.Value.(int64), err)
		}
	}
}
