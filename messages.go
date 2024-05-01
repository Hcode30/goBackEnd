package main

import "time"

type Message struct {
	ID         uint64
	SenderID   uint64
	ReceiverID uint64
	Content    string
	Time       time.Time
}
