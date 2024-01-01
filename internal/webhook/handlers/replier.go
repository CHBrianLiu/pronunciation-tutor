package handlers

import "github.com/line/line-bot-sdk-go/v7/linebot"

type ReplyDoer interface {
	Do() (*linebot.BasicResponse, error)
}

type MessageReplier interface {
	ReplyMessage(replyToken string, messages ...linebot.SendingMessage) ReplyDoer
}
