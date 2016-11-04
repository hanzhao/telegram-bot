package bot

import (
	"log"
	"time"
)

type (
	// Top-level framework instance.
	Bot struct {
		token    string
		handlers []HandlerFunc
	}

	// Bot running mode.
	Mode int

	// HandlerFunc defines a function to resolve updates. Returns true or an
	// error will terminate the handlers chain.
	HandlerFunc func(*Bot, *Update) error
)

func NewBot(token string) *Bot {
	e := &Bot{
		token: token,
	}
	return e
}

func (e *Bot) handle(update *Update) {
	for _, handler := range e.handlers {
		err := handler(e, update)
		if err != nil {
			log.Println("Error:", err, "< handle")
			break
		}
	}
}

func (e *Bot) AddHandler(handler HandlerFunc) {
	e.handlers = append(e.handlers, handler)
}

func (e *Bot) RunWebhook(url string) {
	// TODO
}

func (e *Bot) RunLongPolling() {
	log.Println("Info: Running in long polling mode.")
	offset := 0
	for {
		updates, err := e.GetUpdates(offset, 100, 120)
		if err != nil {
			log.Println("Error:", err, "< GetUpdates < RunLongPolling")
			time.Sleep(time.Second)
			continue
		}
		for _, update := range updates {
			e.handle(&update)
			if offset < (update.UpdateID + 1) {
				offset = update.UpdateID + 1
			}
		}
	}
}
