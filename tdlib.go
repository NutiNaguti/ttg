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

var DefaultTimeout = 1.0

var (
	authorizationState                    = "authorizationState"
	authorizationStateWaitTdlibParameters = "authorizationStateWaitTdlibParameters"
	authorizationStateWaitPhoneNumber     = "authorizationStateWaitPhoneNumber"
	authorizationStateWaitEmailAddress    = "authorizationStateWaitEmailAddress"
	authorizationStateWaitEmailCode       = "authorizationStateWaitEmailCode"
	authorizationStateWaitCode            = "authorizationStateWaitCode"
	authorizationStateWaitRegistration    = "authorizationStateWaitRegistration"
	authorizationStateWaitPassword        = "authorizationStateWaitPassword"
	authorizationStateReady               = "authorizationStateReady"

	setTdlibParameters           = "setTdlibParameters"
	setAuthenticationPhoneNumber = "setAuthenticationPhoneNumber"

	getChats = "getChats"
)

type TdSender interface {
	string | UpdateData
}

type Client struct {
	Client  C.int
	Config  Config
	Updates chan []UpdateData
}

type UpdateData map[string]interface{}

func NewClient() (Client, error) {
	config := GetConfig()
	client := Client{Client: C.td_create_client_id(), Config: *config}
	err := client.setTDLibParams()
	if err != nil {
		log.Error(err)
	}
	return client, nil
}

func (c *Client) AuthorizeViaPhoneNumber(phoneNumber string) error {
	query := UpdateData{
		"@type":        setAuthenticationPhoneNumber,
		"phone_number": phoneNumber,
	}
	err := c.send(query)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (c *Client) GetAllChats() {
	query := UpdateData{
		"@type": getChats,
	}
	c.send(query)
}

func (c *Client) getUpdates(updates chan UpdateData) error {
	err := c.receive(updates, DefaultTimeout)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
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

func (c *Client) receive(updates chan UpdateData, timeout float64) error {
	update := C.td_receive(C.double(timeout))
	var result UpdateData
	log.Info(C.GoString(update))
	if C.GoString(update) == "" {
		return nil
	}
	err := json.Unmarshal([]byte(C.GoString(update)), &result)
	if err != nil {
		log.Error(err)
	}
	updates <- result
	return nil
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
