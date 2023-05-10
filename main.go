package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	botToken := os.Getenv("CHAT_BOT_TOKEN")
	botAPI := "https://api.telegram.org/bot"
	botURL := botAPI + botToken

	infoLogFile, _ := os.Create("/var/log/bot/bot_info.log")
	errorLogFile, _ := os.Create("/var/log/bot/bot_error.log")
	infoLog := log.New(infoLogFile, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(errorLogFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbUser := os.Getenv("TARANTOOL_USER_NAME")
	dbUserPassword := os.Getenv("TARANTOOL_USER_PASSWORD")
	connection := Connection("tarantool_host:3301", dbUser, dbUserPassword, errorLog)
	defer ConnectionClose(connection, errorLog)

	offset := &atomic.Uint64{}
	wg := sync.WaitGroup{}
	done := make(chan bool)

	wg.Add(1)
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				updates := getUpdates(offset, botURL, errorLog)

				for _, update := range updates {
					offset.Store(uint64(update.UpdateId) + 1)

					switch update.Message.Text {
					case "/start":
						infoLog.Println(fmt.Sprintf("User with username - %s joined bot", update.Message.UserInfo.Username))
						_ = sendMessage("Welcome to password manager bot\n\n"+
							"For security reasons, messages containing credits "+
							"will expire in 5 minutes", update, botURL, errorLog)
						createSpaceForUser(update.Message.Chat.ChatId, connection, infoLog, errorLog)
						connection = Connection("tarantool_host:3301", dbUser, dbUserPassword, errorLog)
					case "/set":
						_ = sendMessage("To save the password, send a message according to this format: "+
							"<Service name> <Login> <Password>\n\n"+
							"Type /exit to leave \"set\" menu", update, botURL, errorLog)
						getPasswordFromPerson(botURL, offset, update, connection, errorLog)
					case "/get":
						_ = sendMessage("To get a password, send me a service name\n\n"+
							"Type /exit to leave \"get\" menu", update, botURL, errorLog)
						sendPasswordToPerson(botURL, offset, update, connection, errorLog)
					case "/del":
						_ = sendMessage("To delete a password, send me a service name\n\n"+
							"Type /exit to leave \"del\" menu", update, botURL, errorLog)
						deletePassword(botURL, offset, update, connection, errorLog)
					}
				}
			}
		}
	}()
	wg.Wait()
}
