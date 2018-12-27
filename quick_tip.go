package main

import (
	"math/rand"
	"time"
)

var messages = [...]string{"Drink water many times a day", "Stretch every morning", "Exercise at least 3 times a week"}

func say() string {
	rand.Seed(time.Now().Unix())
	return messages[rand.Intn(len(messages))]
}
