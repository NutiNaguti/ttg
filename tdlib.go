package main

/*
#cgo CFLAGS: -I./td/include
#cgo LDFLAGS: -L./td/lib -ltdjson
#include <stdlib.h>
#include <string.h>
#include <td/telegram/td_json_client.h>
*/
import "C"
import (
	"encoding/json"
	"unsafe"

	"github.com/charmbracelet/log"
)

var DefaultTimeout = 3.0

var (
	authorizationStateWaitTdlibParameters = "authorizationStateWaitTdlibParameters"
	authorizationStateWaitPhoneNumber     = "authorizationStateWaitPhoneNumber"
	authorizationStateWaitEmailAddress    = "authorizationStateWaitEmailAddress"
	authorizationStateWaitEmailCode       = "authorizationStateWaitEmailCode"
	authorizationStateWaitCode            = "authorizationStateWaitCode"
	authorizationStateWaitRegistration    = "authorizationStateWaitRegistration"
	authorizationStateWaitPassword        = "authorizationStateWaitPassword"
	authorizationStateReady               = "authorizationStateReady"

	updateOption = "updateOption"

	setTdlibParameters = "setTdlibParameters"
)

type TdSender interface {
	string | UpdateData
}

type Client struct {
	Client/* unsafe.Pointer */ C.int
	Config Config
}

type UpdateData map[string]interface{}

func NewClient(c chan *Client) error {
	config := GetConfig()
	client := Client{Client: C.td_create_client_id(), Config: *config}
	client.setTDLibParams()
	c <- &client
	return nil
}

func (c *Client) GetAllChats() {}

func (c *Client) getUpdates() (UpdateData, error) {
	updates, err := c.receive()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return updates, nil
}

func (c *Client) setTDLibParams() error {
	params := UpdateData{
		"@type":                setTdlibParameters,
		"api_id":               c.Config.ApiId,
		"api_hash":             c.Config.ApiHash,
		"database_directory":   c.Config.DatabaseDirectory,
		"use_message_database": c.Config.UseMessageDatabase,
		"use_secret_chats":     c.Config.UseSecretChats,
		"system_language_code": c.Config.SystemLanguageCode,
		"device_model":         c.Config.DeviceModel,
		"application_version":  c.Config.AppVersion,
	}

	err := c.send(params)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (c *Client) receive() (UpdateData, error) {
	update := C.td_receive(C.double(1))
	var result UpdateData
	log.Info(C.GoString(update))
	err := json.Unmarshal([]byte(C.GoString(update)), &result)
	if err != nil {
		log.Error(err)
	}
	return result, nil
}

func (c *Client) send(jsonQuery UpdateData) error {
	var query *C.char

	jsonBytes, err := json.Marshal(jsonQuery)
	if err != nil {
		log.Error(err)
		return err
	}

	query = (*C.char)(C.CString(string(jsonBytes)))
	defer C.free(unsafe.Pointer(query))

	log.Info(C.GoString(query))

	C.td_send(c.Client, query)
	return nil
}
