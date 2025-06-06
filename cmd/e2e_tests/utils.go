package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// NewRandomUserName generates a random name suffixed with current time + offset.
func NewRandomUserName(prefix string, nameLen int, timeSuffixOffset time.Duration) string {
	_unixMilli := time.Now().Add(timeSuffixOffset).UnixMilli()
	unixMilli := strconv.Itoa(int(_unixMilli))
	b := make([]rune, nameLen)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return fmt.Sprintf("%s_%s_%s", prefix, string(b), unixMilli)
}
