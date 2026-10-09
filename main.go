package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Activity struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"repo"`
	Payload struct {
		Action  string            `json:"action"`
		Commits []json.RawMessage `json:"commits"`
	} `json:"payload"`
}

func fetchEvents(username string) ([]Activity, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Git-Tracker-app")

	token := os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer"+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Не удалось подключиться к github: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error status code: %d", resp.StatusCode)
	}

	var activity []Activity
	if err := json.NewDecoder(resp.Body).Decode(&activity); err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}

	return activity, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Не удалось загрузить .env: ", err)
		return
	}

	username := os.Args[1]

	git_data, err := fetchEvents(username)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if git_data == nil {
		fmt.Println("Человек лох")
		return
	}

	for i, event := range git_data {
		fmt.Printf("%d. Type: %s\n", i+1, event.Type)
		fmt.Printf("   Repo: %s\n", event.Repo.Name)
		fmt.Printf("   URL: %s\n", event.Repo.Url)
		fmt.Printf("   Action: %s\n\n", event.Payload.Action)
	}
}
