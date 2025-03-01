package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
    // Start
	log.Println("Starting the program...")

	// Required environment variables
	webhookURL := os.Getenv("HNOTIFY_DISCORD_WEBHOOK_URL")
	username := os.Getenv("HNOTIFY_DISCORD_USERNAME")
	jsonURL := os.Getenv("HNOTIFY_JSON_URL")
	localFilename := os.Getenv("HNOTIFY_LOCAL_FILENAME")

	// Check if required environment variables are set
	if webhookURL == "" {
		log.Fatal("Error: HNOTIFY_DISCORD_WEBHOOK_URL environment variable is not set.")
	}
	if jsonURL == "" {
		jsonURL = "https://github.com/arkadiyt/bounty-targets-data/raw/refs/heads/main/data/hackerone_data.json"
	}

	if username == "" {
		username = "Notification Bot"
	}
	if localFilename == "" {
		localFilename = "data.json"
	}

	// Print the values to verify
	fmt.Println("Username:", username)
	fmt.Println("Webhook URL:", webhookURL)
	fmt.Println("JSON URL:", jsonURL)
	fmt.Println("Local Filename:", localFilename)

	// Download the latest JSON data
	newData, err := downloadJSON(jsonURL)
	if err != nil {
		log.Println("Error downloading JSON:", err)
		return
	}

	// If the local file doesn't exist, save the initial JSON and exit
	if !fileExists(localFilename) {
		log.Println("Local file does not exist. Saving initial JSON...")
		if err := saveJSONToFile(newData, localFilename); err != nil {
			log.Println("Error saving initial JSON:", err)
			return
		}
		log.Println("Initial JSON downloaded and saved.")
		return
	}

	// Read the local JSON file
	oldData, err := readJSONFromFile(localFilename)
	if err != nil {
		log.Println("Error reading local file:", err)
		return
	}

	// Process the old and new JSON data
	process(oldData, newData, username, localFilename, webhookURL)

	log.Println("Program execution completed.")
}

// process compares old and new JSON data and sends notifications for changes.
func process(oldData, newData []byte, username, localFilename, webhookURL string) {
	log.Println("Processing JSON data...")
	changes, err := compareJSON(oldData, newData)
	if err != nil {
		log.Println("Error comparing JSON:", err)
		return
	}
	if changesLen := len(changes); changesLen > 0 {
		log.Printf("Found %d changes in JSON data.\n", changesLen)

		log.Println("Processing notifications...")
		// Send a notification for each change
		var httpClient = createHTTPClient()
		for _, change := range changes {
			sendNotification(httpClient, webhookURL, username, change)
		}
		log.Println("All Notifications sent.")

		// Save the new JSON data as the new local copy
		if err := saveJSONToFile(newData, localFilename); err != nil {
			log.Println("Error saving new JSON data:", err)
			return
		}
	} else {
		log.Println("No changes detected.")
	}
}
