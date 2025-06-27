// Package utils provides utility functions for the M-Pesa Daraja API.
package utils

import (
	"encoding/base64"
	"fmt"
	"time"
)

// CheckKeys validates that the required keys are present in the payload.
func CheckKeys(requiredKeys []string, payload map[string]interface{}) (map[string]interface{}, error) {
	cleanedPayload := make(map[string]interface{})
	for _, key := range requiredKeys {
		if value, exists := payload[key]; exists {
			cleanedPayload[key] = value
		} else {
			return nil, fmt.Errorf("missing key: %s", key)
		}
	}
	return cleanedPayload, nil
}

// GetTimestamp returns the current timestamp in YYYYMMDDHHMMSS format.
func GetTimestamp() string {
	return time.Now().Format("20060102150405")
}

// EncodePassword encodes the password using shortcode, passkey, and timestamp.
func EncodePassword(shortcode, passkey string) string {
	timestamp := GetTimestamp()
	data := shortcode + passkey + timestamp
	return base64.StdEncoding.EncodeToString([]byte(data))
}
