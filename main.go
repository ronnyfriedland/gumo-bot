package main

import (
	"bytes"
	"encoding/json"
	"github.com/magiconair/properties"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

const gumoMessages = "gumo.messages"
const gumoStatus = "gumo.status"

func main() {
	userId, authToken, channel, url := readProperties()

	if needToGumo(gumoStatus) {
		data, _ := json.Marshal(map[string]string{
			"channel": channel,
			"text":    chooseMessage(gumoMessages),
		})

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			log.Fatalf("Error creating request %v", err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Auth-Token", authToken)
		req.Header.Add("X-User-Id", userId)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			log.Printf("Successfully sent message for today")
			updateNeedToGumo(gumoStatus)
		} else {
			log.Printf("Error sending message for today %d", resp.StatusCode)
		}
	} else {
		log.Printf("Already sent message for today - skipping")
	}
}

func readProperties() (string, string, string, string) {
	props := properties.MustLoadFile("gumo.properties", properties.UTF8)
	userId := props.GetString("userId", "unknown-userId")
	authToken := props.GetString("authToken", "unknown-authToken")
	channel := props.GetString("channel", "unknown-channel")
	url := props.GetString("url", "unknown-url")
	return userId, authToken, channel, url
}

func needToGumo(filename string) bool {
	currentTime := time.Now()
	statusDateString := currentTime.Format(dateLayout)

	body, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read file %v", err)
	}
	lastStatusDateString := string(body)

	if lastStatusDateString != "" {
		statusDate, _ := time.Parse(time.DateOnly, statusDateString)
		lastStatusDate, _ := time.Parse(time.DateOnly, lastStatusDateString)
		return statusDate.After(lastStatusDate)
	} else {
		return true
	}

}

func updateNeedToGumo(filename string) {
	currentTime := time.Now()
	statusDateString := currentTime.Format(dateLayout)

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("Unable to create status file %v", err)
	}

	defer f.Close()

	_, err2 := f.WriteString(statusDateString)

	if err2 != nil {
		log.Fatalf("Unable to write to status file %v", err2)
	}
}

func chooseMessage(filename string) string {
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

func shuffleMessage(messages []string) string {
	rand.Shuffle(len(messages), func(i, j int) {
		messages[i], messages[j] = messages[j], messages[i]
	})
	return messages[0]
}
