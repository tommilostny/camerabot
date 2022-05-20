package camerabot

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tommilostny/camerabot/telegram"
)

const fallbackTimeout = 20 * time.Second

var (
	// MainChatID refers to chat where camerabot will send error reports.
	MainChatID int64

	// CacheDir is a path to dorectory where last photos are stored.
	CacheDir string

	// Chat IDs set in AllowedChatIds are the only ones camerabot will respond to.
	AllowedChatIDs *list.List

	mu           sync.Mutex
	lastUpdateID int64
)

// Handlers implement commands that are executed by bot. Unknown commands ignored.
var Handlers = make(map[string]Handler)

// Handler processes command sent to bot.
type Handler interface {
	// Command name supported by handler.
	Command() string
	// Handle supported command.
	Handle(chatID int64) error
	// Help message. Help handler will show it.
	Help() string
}

func command(u telegram.Update) string {
	if len(u.Message.Entities) == 0 {
		return ""
	}

	if u.Message.Entities[0].Type == "bot_command" {
		return u.Message.Text
	}

	return ""
}

func trackLastUpdateID(ID int64) {
	mu.Lock()
	log.Printf("Last update ID: %d, incoming update ID: %d", lastUpdateID, ID)
	if lastUpdateID < ID {
		lastUpdateID = ID
	}
	mu.Unlock()
}

func handleUpdates(updates []telegram.Update) {
	for _, u := range updates {
		trackLastUpdateID(u.ID)

		cmd := command(u)
		chatID := u.Message.Chat.ID

		if cmd == "" || !isAllowedChatID(chatID) {
			continue
		}

		if h, exists := Handlers[cmd]; exists {
			h.Handle(chatID)
			continue
		}

		log.Printf("Unknown command: %q in chat: %d ignored", cmd, u.Message.Chat.ID)
	}
}

// ListenAndServe gets updates and processes them.
func ListenAndServe() {
	telegram.SendTextMessage(MainChatID, "Hi there.")

	for {
		updates, err := telegram.GetUpdates(lastUpdateID + 1)
		if err != nil {
			telegram.SendTextMessage(MainChatID, fmt.Sprintf("Failed getting updates: %v", err))
			time.Sleep(fallbackTimeout)
		}

		log.Print("Polling...")
		handleUpdates(updates)
	}
}

func ParseAllowedChatIDs() *list.List {
	splitted := strings.Split(os.Getenv("ALLOWED_CHAT_IDS"), ";")
	if len(splitted) == 0 {
		return nil
	}

	result := list.New()
	for i := 0; i < len(splitted); i++ {

		chatID, err := strconv.ParseInt(splitted[i], 0, 64)
		if err != nil {
			log.Println(err)
			continue
		}

		result.PushBack(chatID)
		log.Println("Allowed chat ID:", chatID)
	}

	if result.Len() == 0 {
		return nil
	}
	return result
}

func isAllowedChatID(chatID int64) bool {
	if AllowedChatIDs == nil {
		return true
	}

	for e := AllowedChatIDs.Front(); e != nil; e = e.Next() {
		if e.Value.(int64) == chatID {
			return true
		}
	}

	log.Println("Chat ID:", chatID, "is not allowed")
	return false
}
