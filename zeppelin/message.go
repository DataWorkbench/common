package zeppelin

import (
	"strings"

	"github.com/valyala/fastjson"
)

type MessageHandler interface {
	onMessage(msg string) error
}

type OutputHandler interface {
	onStatementAppendOutput(statementId string, index int, output string) error
	onStatementUpdateOutput(statementId string, index int, oType string, output string) error
}

type AbstractMessageHandler struct {
	OutputHandler
}

func (handler *AbstractMessageHandler) onMessage(msg string) error {
	msgJson, err := fastjson.Parse(msg)
	if err != nil {
		return err
	}
	op := string(msgJson.GetStringBytes("op"))
	if strings.EqualFold(op, "PARAGRAPH_UPDATE_OUTPUT") {
		paragraphId := string(msgJson.Get("data").GetStringBytes("paragraphId"))
		index := msgJson.GetInt("index")
		mType := string(msgJson.Get("data").GetStringBytes("type"))
		output := string(msgJson.Get("data").GetStringBytes("data"))
		return handler.onStatementUpdateOutput(paragraphId, index, mType, output)
	}
	if strings.EqualFold(op, "PARAGRAPH_APPEND_OUTPUT") {
		paragraphId := string(msgJson.Get("data").GetStringBytes("paragraphId"))
		index := msgJson.GetInt("index")
		output := string(msgJson.Get("data").GetStringBytes("data"))
		return handler.onStatementAppendOutput(paragraphId, index, output)
	}
	return nil
}

type CompositeMessageHandler struct {
	*AbstractMessageHandler
	MessageHandler map[string]OutputHandler
}

func (handler *CompositeMessageHandler) onStatementAppendOutput(statementId string, index int, output string) error {
	messageHandler := handler.MessageHandler[statementId]
	if messageHandler == nil {
		return nil
	}
	return messageHandler.onStatementAppendOutput(statementId, index, output)
}

func (handler *CompositeMessageHandler) onStatementUpdateOutput(statementId string, index int, oType string, output string) error {
	messageHandler := handler.MessageHandler[statementId]
	if messageHandler == nil {
		return nil
	}
	return messageHandler.onStatementUpdateOutput(statementId, index, oType, output)
}

func (handler *CompositeMessageHandler) addStatementMessageHandler(statementId string, outputHandler OutputHandler) {
	handler.MessageHandler[statementId] = outputHandler
}
