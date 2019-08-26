package vk_objects

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/SevereCloud/vksdk/5.92/object"
)

type Message struct {
	UserID     *int                    `json:"user_id,omitempty"`
	UserIDs    *string                 `json:"user_ids,omitempty"`
	Message    string                  `json:"message"`
	RandomId   int32                   `json:"random_id"`
	PeerId     *int                    `json:"peer_id,omitempty"`
	Domain     *string                 `json:"domain,omitempty"`
	ChatId     *int                    `json:"chat_id,omitempty"`
	Latitude   *float64                `json:"lat,omitempty"`
	Longitude  *float64                `json:"long,omitempty"`
	Attachment *string                 `json:"attachment,omitempty"`
	ReplyTo    *int                    `json:"reply_to,omitempty"`
	StickerId  *int                    `json:"sticker_id,omitempty"`
	GroupId    *int                    `json:"group_id,omitempty"`
	Keyboard   object.MessagesKeyboard `json:"keyboard"`
}

func NewMessage() Message {
	message := Message{}

	keyboard := object.MessagesKeyboard{
		OneTime: false,
		Buttons: make([][]object.MessagesKeyboardButton, 0),
	}

	message.Keyboard = keyboard

	return message
}

func (m *Message) ToMap() map[string]string {
	mapped := make(map[string]string, 0)

	iVal := reflect.ValueOf(m).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)

		var v string
		switch f.Interface().(type) {
		case *int, *int8, *int16, *int32, *int64:
			vPtr := f.Elem()
			if !vPtr.IsValid() {
				continue
			}
			v = strconv.FormatInt(f.Elem().Int(), 10)
		case int, int8, int16, int32, int64:
			v = strconv.FormatInt(f.Int(), 10)
		case uint, uint8, uint16, uint32, uint64:
			v = strconv.FormatUint(f.Uint(), 10)
		case float32:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 32)
		case float64:
			v = strconv.FormatFloat(f.Float(), 'f', 4, 64)
		case []byte:
			v = string(f.Bytes())
		case string:
			v = f.String()
		}

		if v != "" {
			jsonTag := typ.Field(i).Tag.Get("json")
			if jsonTag != "" {
				var name string
				splitted := strings.Split(jsonTag, ",")
				if len(splitted) > 1 {
					name = splitted[0]
				} else {
					name = jsonTag
				}
				mapped[name] = v
			}
		}
	}
	if len(m.Keyboard.Buttons) > 0 {
		marshal, _ := json.Marshal(m.Keyboard)
		mapped["keyboard"] = string(marshal)
	}
	return mapped
}
