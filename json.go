package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// downloadJSON downloads JSON data from the specified URL.
func downloadJSON(url string) ([]byte, error) {
	log.Printf("Downloading JSON from %s...\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download JSON: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// saveJSONToFile saves JSON data to a local file.
func saveJSONToFile(data []byte, filename string) error {
	log.Printf("Saving JSON to file: %s\n", filename)
	return os.WriteFile(filename, data, 0644)
}

// readJSONFromFile reads JSON data from a local file.
func readJSONFromFile(filename string) ([]byte, error) {
	log.Printf("Reading JSON from file: %s\n", filename)
	return os.ReadFile(filename)
}

// fileExists checks if a file exists.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// compareJSON compares old and new JSON data and detects changes.
func compareJSON(oldData, newData []byte) ([]string, error) {
	log.Println("Comparing old and new JSON data...")
	var oldPrograms, newPrograms []map[string]interface{}
	if err := json.Unmarshal(oldData, &oldPrograms); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(newData, &newPrograms); err != nil {
		return nil, err
	}

	// Create a map of old programs for quick lookup
	oldProgramMap := make(map[string]map[string]interface{})
	for _, program := range oldPrograms {
		if handle, ok := program["handle"].(string); ok {
			oldProgramMap[handle] = program
		}
	}

	var changes []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Compare new programs with old ones concurrently
	for _, newProgram := range newPrograms {
		wg.Add(1)
		go func(newProgram map[string]interface{}) {
			defer wg.Done()

			handle, ok := newProgram["handle"].(string)
			if !ok {
				return
			}

			oldProgram, exists := oldProgramMap[handle]
			if !exists {
				// New program added
				mu.Lock()
				changes = append(changes, fmt.Sprintf(
					"New PROGRAM Added:\n"+
						"- Name: %s\n"+
						"- Handle: %s\n"+
						"- URL: %s\n"+
						"- Offers Bounties: %v\n"+
						"- Offers Swag: %v\n"+
						"- Response Efficiency: %.0f%%\n"+
						"- Submission State: %s",
					newProgram["name"], handle, newProgram["url"], newProgram["offers_bounties"], newProgram["offers_swag"], newProgram["response_efficiency_percentage"], newProgram["submission_state"],
				))
				mu.Unlock()
				return
			}

			// Safely extract in-scope targets
			var newTargets, oldTargets []interface{}

			if targets, ok := newProgram["targets"].(map[string]interface{}); ok {
				if inScope, ok := targets["in_scope"].([]interface{}); ok {
					newTargets = inScope
				}
			}

			if targets, ok := oldProgram["targets"].(map[string]interface{}); ok {
				if inScope, ok := targets["in_scope"].([]interface{}); ok {
					oldTargets = inScope
				}
			}

			// Create a map of old targets for quick lookup
			oldTargetMap := make(map[string]struct{})
			for _, target := range oldTargets {
				if targetMap, ok := target.(map[string]interface{}); ok {
					if asset, ok := targetMap["asset_identifier"].(string); ok {
						oldTargetMap[asset] = struct{}{}
					}
				}
			}

			// Check for new targets
			for _, target := range newTargets {
				if targetMap, ok := target.(map[string]interface{}); ok {
					if asset, ok := targetMap["asset_identifier"].(string); ok {
						if _, exists := oldTargetMap[asset]; !exists {
							mu.Lock()
							changes = append(changes, fmt.Sprintf(
								"New ASSET Added:\n"+
									"- Program: %s (%s)\n"+
									"- Asset Type: %s\n"+
									"- Asset: %s\n"+
									"- Eligible for Bounty: %v\n"+
									"- Max Severity: %s",
								newProgram["name"], handle, targetMap["asset_type"], targetMap["asset_identifier"], targetMap["eligible_for_bounty"], targetMap["max_severity"],
							))
							mu.Unlock()
						}
					}
				}
			}
		}(newProgram)
	}

	wg.Wait()
	return changes, nil
}
