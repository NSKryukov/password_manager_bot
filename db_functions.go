package main

import (
	"fmt"
	"github.com/tarantool/go-tarantool"
	"log"
	"strconv"
)

func Connection(dbUrl, user, password string, errorLog *log.Logger) *tarantool.Connection {
	conn, err := tarantool.Connect(dbUrl, tarantool.Opts{
		User: user,
		Pass: password,
	})
	if err != nil {
		errorLog.Println("Error connecting to db, error: ", err)
	}
	return conn
}

func ConnectionClose(conn *tarantool.Connection, errorLog *log.Logger) {
	err := conn.Close()
	if err != nil {
		errorLog.Println("Error closing connecting with db, error: ", err)
	}
}

func createSpaceForUser(ChatID int, conn *tarantool.Connection, infoLog, errorLog *log.Logger) {
	space := "user" + strconv.Itoa(ChatID)
	_, err := conn.Call("box.schema.space.create", []interface{}{
		space,
		map[string]bool{"if_not_exists": true}})
	if err != nil {
		errorLog.Println("Error while creating space for user, error: ", err)
	}
	infoLog.Println(fmt.Sprintf("Space for user with chat id - %v created", ChatID))

	_, err = conn.Call(fmt.Sprintf("box.space.%s:format", space), [][]map[string]string{
		{
			{"name": "service", "type": "string"},
			{"name": "login", "type": "string"},
			{"name": "password", "type": "string"},
		}})
	if err != nil {
		errorLog.Println("Error while formatting space for user, error: ", err)
	}

	_, err = conn.Call(fmt.Sprintf("box.space.%s:create_index", space), []interface{}{
		"primary",
		map[string]interface{}{
			"parts":         []string{"service"},
			"if_not_exists": true}})
	if err != nil {
		errorLog.Println("Error while indexing user space, error: ", err)
	}
}

func insertData(ChatID int, serviceName, login, password string, conn *tarantool.Connection) (bool, error) {
	space := "user" + strconv.Itoa(ChatID)
	resp, err := conn.Insert(space, []interface{}{serviceName, login, password})
	if resp == nil {
		return true, err
	} else if resp.Code == 3 {
		return false, err
	} else if err != nil {
		return true, err
	}
	return true, err
}

func selectData(ChatID int, serviceName string, conn *tarantool.Connection) (*tarantool.Response, error) {
	space := "user" + strconv.Itoa(ChatID)
	resp, err := conn.Select(space, "primary", 0, 1, tarantool.IterEq, []interface{}{serviceName})
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func deleteData(ChatID int, serviceName string, conn *tarantool.Connection) (*tarantool.Response, error) {
	space := "user" + strconv.Itoa(ChatID)
	resp, err := conn.Delete(space, "primary", []interface{}{serviceName})
	return resp, err
}
