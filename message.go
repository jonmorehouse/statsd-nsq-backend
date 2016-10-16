package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	ErrInvalidMsg     = errors.New("E_INVALID_MSG")
	ErrUnknownMsgType = errors.New("E_UNKNOWN_MSG_TYPE")
)

type Message interface {
	Prefix() string
	Key() string
	Type() string
	Value() interface{}
	Tags() []string
	SampleRate() float64

	RawBytes() []byte
}

func ParseMessage(byts []byte) (Message, error) {
	str := string(byts)
	pieces := strings.Split(str, "|", -1)

	if len(pieces) < 2 {
		return nil, ErrInvalidMsg
	}

	msg := &message{
		rawMsg: str,
	}

	var value interface{}
	var err error

	switch pieces[1] {
	case "c":
		msg.typ = "count"
		value, err = strconv.Atoi(pieces[0])
	case "ms":
		msg.typ = "timing"
		value, err = strconv.ParseFloat(pieces[1])
	case "h":
		msg.typ = "histogram"
		value, err = strconv.ParseFloat(pieces[1])

	case "g":
		msg.typ = "gauge"
		value, err = strconv.ParseFloat(pieces[1])
	case "s":
		msg.typ = "set"
		value = pieces[1]
	default:
		return nil, ErrUnknownMsgType
	}

	// check for a prefix
	keyPieces := strings.Split(pieces[0], ".", 1)
	if len(pieces) > 1 {
		msg.prefix = keyPieces[0]
		msg.key = keyPieces[1]
	} else {
		msg.prefix = ""
		msg.key = keyPieces[0]
	}

	//
	for _, piece := range pieces {

	}

}

type message struct {
	prefix, key, typ, rawMsg string
	value                    interface{}

	timestamp  time.Time
	sampleRate float64
	tags       []string
}

func (m *message) Prefix() string {
	return m.prefix
}

func (m *message) Key() string {
	return m.key
}

func (m *message) Type() string {
	return m.typ
}

func (m *message) Value() interface{} {
	return m.value
}

func (m *message) Tags() []string {
	return m.tags
}

func (m *message) RawMessage() string {
	return m.rawMsg
}
