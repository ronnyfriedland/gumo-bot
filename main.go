package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"ronnyfriedland/gumo/configuration"
	"ronnyfriedland/gumo/message"
	"time"
)

const dateLayout = "2006-01-02"


// THe main method to start the application
func main() {
	var configpath = flag.String("configpath", "/etc/gumo", "the config path")
	var directmessage = flag.String("message", "", "the message to send directly")
	flag.Parse()

	var gumoProperties = *configpath + "/gumo.properties"
	var gumoMessages = *configpath + "/gumo.messages"
	var gumoStatus = *configpath + "/gumo.status"

	if needToGumo(gumoStatus) {
		target, err := configuration.GetTarget(gumoProperties)
		if err != nil {
			log.Fatalf("configuration error: %v", err)
			return
		}

		url, err := target.Url(gumoProperties)
		if err != nil {
			log.Fatalf("Error getting url: %v", err)
			return
		}

		if *directmessage == "" {
			*directmessage = message.ChooseMessage(gumoMessages)
		}

		data, err := target.Payload(gumoProperties, *directmessage)
		if err != nil {
			log.Fatalf("Error building payload: %v", err)
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Error formatting json: %v", err)
			return
		}

		headers, err := target.Headers(gumoProperties)
		if err != nil {
			log.Fatalf("Error getting headers: %v", err)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatalf("Error creating request: %v", err)
			return
		}

		for k := range headers {
			req.Header.Add(k, headers[k])
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Error sending request %v", err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatalf("Unable to close response stream %v", err)
			}
		}(resp.Body)

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

// Check if gumo was already triggered
func needToGumo(filename string) bool {
	currentTime := time.Now()
	statusDateString := currentTime.Format(dateLayout)

	body, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Unable to read file %v", err)
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

// Update the gumo status - set the current date to prevent another gumo today
func updateNeedToGumo(filename string) {
	currentTime := time.Now()
	statusDateString := currentTime.Format(dateLayout)

	f, err := os.Create(filename)
	if err != nil {

		log.Fatalf("Unable to create status file %v", err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Unable to close file %v", err)
		}
	}(f)

	_, err2 := f.WriteString(statusDateString)

	if err2 != nil {
		log.Fatalf("Unable to write to status file %v", err2)
	}
}
