package main

/*
#cgo CFLAGS: -I./td/include
#cgo LDFLAGS: -L./td/lib -ltdjson
#include <td/telegram/td_json_client.h>
*/
import "C"
import (
	"encoding/json"
	"unsafe"
)

var DefaultTimeout = 10.0

type Client struct {
	Client unsafe.Pointer
	Config Config
}

type Config struct {
	ApiId              int32  `json:"api_id"`
	ApiHash            string `json:"api_hash"`
	DatabaseDirectory  string `json:"database_directory"`
	UseMessageDatabase bool   `json:"use_message_database"`
	UseSecretChat      bool   `json:"use_secret_chats"`
	SystemLanguageCode string `json:"system_language_code"`
	DeviceModel        string `json:"device_model"`
}

func NewClient() *Client {
	config := Config{}
	client := Client{Client: C.td_json_client_create(), Config: config}
	return &client
}

func (c *Client) Receive() ([]byte, error) {
	result := C.td_json_client_receive(c.Client, C.double(DefaultTimeout))
	return []byte(C.GoString(result)), nil
}

func (c *Client) Send(jsonQuery interface{}) {
	var data = new(map[string]interface{})
	json.Unmarshal([]byte(jsonQuery.(string)), data)
}
