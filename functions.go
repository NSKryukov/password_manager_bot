package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tarantool/go-tarantool"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func getUpdates(offset *atomic.Uint64, botUrl string, errorLog *log.Logger) []Update {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.FormatUint(offset.Load(), 10))
	if err != nil {
		errorLog.Println("Error while getting updates from bot API, error: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorLog.Println("Error while serializing response body, error: ", err)
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		errorLog.Println("Error while unmarshalling response body, error: ", err)
	}
	return restResponse.Result
}

func getUpdateByID(updates []Update, update Update) Update {
	for _, elem := range updates {
		if elem.Message.Chat.ChatId == update.Message.Chat.ChatId {
			return elem
		}
	}
	return Update{}
}

func sendMessage(text string, Update Update, botUrl string, errorLog *log.Logger) []int {
	var MessageToSend MessageToSend
	MessageToSend.Text = text
	MessageToSend.ChatId = Update.Message.Chat.ChatId

	buf, err := json.Marshal(MessageToSend)
	if err != nil {
		errorLog.Println("Error while marshalling request, error: ", err)
		return nil
	}
	resp, err := http.Post(
		botUrl+"/sendMessage",
		"application/json",
		bytes.NewBuffer(buf),
	)
	if err != nil {
		errorLog.Println("Error sending message through API, error: ", err)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errorLog.Println("Error while serializing response body, error: ", err)
		return nil
	}
	defer resp.Body.Close()

	var ReceivedMessage ReceivedMessage
	err = json.Unmarshal(body, &ReceivedMessage)
	if err != nil {
		errorLog.Println("Error while unmarshalling response body, error: ", err)
	}
	result := []int{ReceivedMessage.Result.Chat.ChatId, ReceivedMessage.Result.MessageID}
	return result
}

func deleteMessage(chatID, messageID int, botUrl string, errorLog *log.Logger) {
	var DeletingMessage MessageToDelete
	DeletingMessage.ChatID = chatID
	DeletingMessage.MessageID = messageID

	buf, err := json.Marshal(DeletingMessage)
	if err != nil {
		errorLog.Println("Error while marshalling request body, error: ", err)
	}
	time.Sleep(time.Minute * 5)
	_, er := http.Post(
		botUrl+"/deleteMessage",
		"application/json",
		bytes.NewBuffer(buf),
	)
	if er != nil {
		errorLog.Println("Error deleting message through API, error: ", err)
	}
}
func getPasswordFromPerson(botURL string, offset *atomic.Uint64, update Update,
	connection *tarantool.Connection, errorLog *log.Logger) {
	for true {
		updates := getUpdates(offset, botURL, errorLog)

		dataUpdate := getUpdateByID(updates, update)
		if dataUpdate != (Update{}) {
			offset.Add(1)
			data := strings.Split(dataUpdate.Message.Text, " ")
			if data[0] == "/exit" {
				_ = sendMessage("Exited context menu", update, botURL, errorLog)
				break
			} else if len(data) == 3 {
				trigger, err := insertData(update.Message.Chat.ChatId, data[0], data[1], data[2], connection)
				if err != nil && trigger {
					errorLog.Println("Error while inserting data to db, error: ", err)
					_ = sendMessage("Error from bot, try again", update, botURL, errorLog)
					break
				}
				if trigger {
					_ = sendMessage("Password added", update, botURL, errorLog)
					go deleteMessage(dataUpdate.Message.Chat.ChatId, dataUpdate.Message.MessageID, botURL, errorLog)
					break
				} else {
					_ = sendMessage("Service password already exist, try again", update, botURL, errorLog)
				}
			} else {
				_ = sendMessage("Wrong message format try again", update, botURL, errorLog)
			}
			go deleteMessage(dataUpdate.Message.Chat.ChatId, dataUpdate.Message.MessageID, botURL, errorLog)
		}
	}
	time.Sleep(time.Millisecond * 500)
}

func sendPasswordToPerson(botURL string, offset *atomic.Uint64, update Update,
	connection *tarantool.Connection, errorLog *log.Logger) {
	for true {
		updates := getUpdates(offset, botURL, errorLog)

		dataUpdate := getUpdateByID(updates, update)
		if dataUpdate != (Update{}) {
			offset.Add(1)
			data := dataUpdate.Message.Text
			if data == "/exit" {
				_ = sendMessage("Exited context menu", update, botURL, errorLog)
				break
			} else if strings.ContainsAny(data, " ") {
				_ = sendMessage("Wrong message format try again", update, botURL, errorLog)
			} else {
				resp, err := selectData(dataUpdate.Message.Chat.ChatId, data, connection)
				if err != nil {
					errorLog.Println("Can't select data from space, error: ", err)
				}
				selectedTuple := getTupleFromTarantoolResponse(resp)
				if len(selectedTuple) == 3 {
					messageInfo := sendMessage(fmt.Sprintf("Service: %s\nLogin: %s\npassword: %s",
						selectedTuple[0], selectedTuple[1], selectedTuple[2]), update, botURL, errorLog)
					go deleteMessage(messageInfo[0], messageInfo[1], botURL, errorLog)
					break
				} else {
					_ = sendMessage("Service password does not exist, make sure that service name is "+
						"written correctly and try again", update, botURL, errorLog)
				}
			}

		}
	}
	time.Sleep(time.Millisecond * 500)
}

func deletePassword(botURL string, offset *atomic.Uint64, update Update, connection *tarantool.Connection, errorLog *log.Logger) {
	for true {
		updates := getUpdates(offset, botURL, errorLog)

		dataUpdate := getUpdateByID(updates, update)
		if dataUpdate != (Update{}) {
			offset.Add(1)
			data := dataUpdate.Message.Text
			if data == "/exit" {
				_ = sendMessage("Exited context menu", update, botURL, errorLog)
				break
			} else if strings.ContainsAny(data, " ") {
				_ = sendMessage("Wrong message format try again", update, botURL, errorLog)
			} else {
				resp, err := deleteData(update.Message.Chat.ChatId, data, connection)
				if err != nil {
					errorLog.Println("Error while deleting data from db, error: ", err)
				}
				deletedTuple := getTupleFromTarantoolResponse(resp)
				if len(deletedTuple) == 0 {
					_ = sendMessage("Service password does not exist, make sure that the service is "+
						"written correctly and try again", update, botURL, errorLog)
				} else {
					_ = sendMessage(fmt.Sprintf("Password from service '%s' deleted",
						deletedTuple[0]), update, botURL, errorLog)
					break
				}
			}
		}
	}
	time.Sleep(time.Millisecond * 500)
}

func getTupleFromTarantoolResponse(response *tarantool.Response) []string {
	var result []string
	for _, row := range response.Tuples() {
		for _, tuple := range row {
			result = append(result, fmt.Sprintf("%v", tuple))
		}
	}
	return result
}
