package message

import (
	"log"
	"math/rand"
	"os"
	"strings"
)

// ChooseMessage Select the gumo message out of the available message
// in message file provided by the given parameter "filename"
func ChooseMessage(filename string) string {
	var messages []string

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read from messages file %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		messages = append(messages, line)
	}
	return shuffleMessage(messages)
}

// Select random message from list of messages
func shuffleMessage(messages []string) string {
	rand.Shuffle(len(messages), func(i, j int) {
		messages[i], messages[j] = messages[j], messages[i]
	})
	return messages[0]
}
