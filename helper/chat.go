package helper

import tdlib "github.com/c0re100/gotdlib/client"

type CommandMessage struct {
	ChatId      int64
	SendToId    int64
	ThreadId    int64
	MsgId       int64
	ReplyTo     int64
	UserId      int64
	QueryId     tdlib.JsonInt64
	Payload     string
	Command     string
	Arg         string
	Text        string
	NoFwd       bool
	Spoiler     bool
	IgnoreCache bool
}
