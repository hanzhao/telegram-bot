Telegram Bot Framework for Go
===

# Install
```
go get -u github.com/magicae/telegram-bot
```
# Usage
```go
import (
	. "github.com/magicae/telegram-bot"
)

func helloWorldHandler(e *Bot, update *Update) error {
	if update.Message != nil {
		_, err := e.SendMessage(&SendMessageRequest{
			ChatID: update.Message.Chat.ID,
			Text: "Hello world",
		})
		return err
	}
	return nil
}

func main() {
	e := NewBot("YOUR_BOT_TOKEN_HERE")
	me, err := e.GetMe()
	if err != nil {
		panic("Error: " + err.Error())
	} else {
		log.Println("Bot info:", me)
	}
	e.AddHandler(helloWorldHandler)
	e.RunLongPolling()
}
```
