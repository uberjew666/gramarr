package main

import (
tb "gopkg.in/telebot.v3"
	"strings"
)

func (e *Env) HandleFallback(m *tb.Message) {
	var msg []string
	msg = append(msg, "I'm sorry, I don't recognize that command.")
	msg = append(msg, "Type /help to see the available bot commands.")
	Send(e.Bot, m.Sender, strings.Join(msg, "\n"))
}
