package controllers

import (
	// "container/list"
	// "github.com/astaxie/beego"
	// "github.com/gorilla/websocket"
	// "time"
	// "encoding/json"
	// "strings"
	"errors"
	"fmt"
)

type mappedValue map[string]interface{}

func (m mappedValue) Getter(keys ...string) (values []interface{}, err error) {
	for _, key := range keys {
		if value, ok := m[key]; ok {
			values = append(values, value)
		} else {
			err = errors.New(fmt.Sprintf("key %s 不存在", key))
			return
		}
	}
	return
}
