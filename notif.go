package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

// NotifData_POST structure for Discord webhook post data
type NotifData_POST struct {
	Content   string `json:"content"`
	Username  string `json:"username,omitempty"`   // Optional: Set a custom username
	AvatarURL string `json:"avatar_url,omitempty"` // Optional: Set a custom avatar URL
}

// wrapLinks prevents Discord from generating previews by enclosing links in angle brackets
func wrapLinks(text string) string {
	re := regexp.MustCompile(`(https?://\S+)`)
	return re.ReplaceAllString(text, `<$1>`)
}

// sendNotification sends a notification to the Discord webhook.
func sendNotification(httpClient *http.Client, username, message, webhookURL string) {
	log.Printf("Sending notification:\n%s\n", message)
	webhook := NotifData_POST{
		Content:  wrapLinks(message),
		Username: username,
	}

	// Convert the webhook struct to JSON
	jsonData, err := json.Marshal(webhook)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	for {
		// Send the POST request to the webhook URL
		resp, err := httpClient.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Error sending notification:", err)
			return
		}
		defer resp.Body.Close()

		// Check if the response is successful
		if resp.StatusCode == http.StatusNoContent {
			log.Println("Notification sent successfully!")
			return
		} else if resp.StatusCode == http.StatusTooManyRequests {
			var rateLimitResp struct {
				RetryAfter float64 `json:"retry_after"` // in seconds, often fractional
			}
			body, _ := io.ReadAll(resp.Body)
			if err = json.Unmarshal(body, &rateLimitResp); err != nil {
				log.Println("Rate limited, but couldn't parse retry_after. Sleeping 2s as fallback.")
				time.Sleep(2 * time.Second)
			} else {
				log.Printf("Rate limited. Retrying after %.2f seconds...\n", rateLimitResp.RetryAfter)
				wait := time.Duration(rateLimitResp.RetryAfter * float64(time.Second))
				time.Sleep(wait)
			}
			continue
		} else {
			log.Printf("Failed to send notification. Status code: %d\n", resp.StatusCode)
			return
		}
	}
}
