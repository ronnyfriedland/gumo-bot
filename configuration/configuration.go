package configuration

import (
	"errors"
	"github.com/magiconair/properties"
)

// Target The target interface
type Target interface {
	Headers(string) (map[string]string, error)
	Payload(string, string) (map[string]string, error)
	Url(string) (string, error)
}

// RocketchatTarget The rocketchat target
type RocketchatTarget struct {
}

// WebexTarget the webex target
type WebexTarget struct {
}

// Url Returns the rocketchat url
func (rc RocketchatTarget) Url(properties string) (string, error) {
	url, err := readProperty(properties, "url")
	if err == nil {
		return url, nil
	}
	return "", err
}

// Payload Returns the rocketchat related payload data
func (rc RocketchatTarget) Payload(properties string, message string) (map[string]string, error) {
	channel, err := readProperty(properties, "channel")
	if err == nil {
		return map[string]string{
			"channel": channel,
			"text":    message,
		}, nil
	}
	return nil, err
}

// Headers Returns the rocketchat related request headers
func (rc RocketchatTarget) Headers(properties string) (map[string]string, error) {
	userId, err := readProperty(properties, "userId")
	if err != nil {
		return nil, err
	}
	authToken, err := readProperty(properties, "authToken")
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	headers["X-User-Id"] = userId
	headers["X-Auth-Token"] = authToken
	headers["Content-Type"] = "application/json"

	return headers, nil
}

// Url Returns the webex url
func (we WebexTarget) Url(properties string) (string, error) {
	url, err := readProperty(properties, "url")
	if err == nil {
		return url, nil
	}
	return "", err
}

// Payload Returns the webex related payload data
func (we WebexTarget) Payload(properties string, message string) (map[string]string, error) {
	roomId, err := readProperty(properties, "channel")
	if err == nil {
		return map[string]string{
			"roomId": roomId,
			"text":   message,
		}, nil
	}
	return nil, err
}

// Headers Returns the webex related request headers
func (we WebexTarget) Headers(properties string) (map[string]string, error) {
	authToken, err := readProperty(properties, "authToken")
	if err != nil {
		return nil, err
	}

	headers := make(map[string]string)
	headers["X-Auth-Token"] = authToken
	headers["Content-Type"] = "application/json"

	return headers, nil
}

// Read property from property file, if not set error is raised
func readProperty(filename string, key string) (string, error) {
	value := readPropertyDefault(filename, key, "")

	if value == "" {
		return "", errors.New("no property with given key '" + key + "' found")
	} else {
		return value, nil
	}
}

// Read property from property file using default value if not set
func readPropertyDefault(filename string, key string, defaultValue string) string {
	props := properties.MustLoadFile(filename, properties.UTF8)
	return props.GetString(key, defaultValue)
}

// GetTarget Returns which target tool hsa to be used
func GetTarget(filename string) (Target, error) {
	target, err := readProperty(filename, "target")

	if err == nil {
		var t Target
		if target == "rocketchat" {
			t = RocketchatTarget{}
		} else if target == "webex" {
			t = WebexTarget{}
		}
		return t, nil
	} else {
		return nil, err
	}
}
