package main

import "errors"

var (
	ErrInvalidMsg     = errors.New("E_INVALID_MSG")
	ErrUnknownMsgType = errors.New("E_UNKNOWN_MSG_TYPE")

	ErrPublishJSON   = errors.New("E_PUBLISH_JSON_MARSHAL")
	ErrPublishBuffer = errors.New("E_PUBLISHER_BUFFER_WRITE")
)
