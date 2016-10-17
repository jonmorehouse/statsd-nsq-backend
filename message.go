package main

import (
	"strconv"
	"strings"
)

type Message struct {
	Prefix     string      `json:"prefix"`
	Key        string      `json:"key"`
	Type       string      `json:"type"`
	Tags       []string    `json:"tags"`
	SampleRate float64     `json:"sample_rate"`
	Value      interface{} `json:"value"`
	RawMsg     string      `json:"raw_message"`
}

func ParseMessage(byts []byte) (*Message, error) {
	str := string(byts)
	pieces := strings.Split(str, "|")

	if len(pieces) < 2 {
		return nil, ErrInvalidMsg
	}

	msg := &Message{
		RawMsg: str,
	}

	var value interface{}
	var err error

	switch pieces[1] {
	case "c":
		msg.Type = "count"
		value, err = strconv.Atoi(pieces[0])
	case "ms":
		msg.Type = "timing"
		value, err = strconv.ParseFloat(pieces[1], 64)
	case "h":
		msg.Type = "histogram"
		value, err = strconv.ParseFloat(pieces[1], 64)

	case "g":
		msg.Type = "gauge"
		value, err = strconv.ParseFloat(pieces[1], 64)
	case "s":
		msg.Type = "set"
		value = pieces[1]
	default:
		return nil, ErrUnknownMsgType
	}

	if err != nil {
		return nil, ErrInvalidMsg
	}

	msg.Value = value

	// check for a key prefix / namespace
	keyPieces := strings.Split(pieces[0], ".")
	if len(pieces) > 1 {
		msg.Prefix = keyPieces[0]
		msg.Key = keyPieces[1]
	} else {
		msg.Prefix = ""
		msg.Key = keyPieces[0]
	}

	// for each piece available, check for additional, known datadog tags such as `|#` or `|h` (tags and hostname respectively)
	for _, piece := range pieces {
		// handle data dog tags
		if strings.HasPrefix(piece, "|#") {
			rawTags := strings.Split(piece, "|#")[1]
			msg.Tags = strings.Split(rawTags, ",")
		} else if strings.HasPrefix(piece, "|@") {
			// handle encoded sample rate
			rawSampleRate := strings.Split(piece, "|@")[1]
			sampleRate, err := strconv.ParseFloat(rawSampleRate, 64)
			if err != nil {
				return nil, ErrInvalidMsg
			}
			msg.SampleRate = sampleRate
		}
	}

	return msg, nil
}
