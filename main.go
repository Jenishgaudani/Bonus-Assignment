package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	//"net/url"
	"os"

	"github.com/dghubble/oauth1"
)

const (
	apiV2Url       = "https://api.x.com/2/tweets"
	deleteTweetUrl = "https://api.x.com/2/tweets/%s"

// const apiV1Url = "https://api.x.com/1.1/statuses/update.json"
)

// Struct to define the JSON payload for API v2
// type TweetPayload struct {
// Text string `json:"text"`
// }
func postTweet(status string) (string, error) {
	apiKey := os.Getenv("TWITTER_API_KEY")
	apiSecretKey := os.Getenv("TWITTER_API_SECRET_KEY")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	
	config := oauth1.NewConfig(apiKey, apiSecretKey)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	// HTTP client with OAuth1 credentials
	httpClient := config.Client(oauth1.NoContext, token)
	payload := map[string]string{
		"text": status,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %v", err)
	}
	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", apiV2Url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	// Set the required headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	// Create an HTTP client and send the request
	//client := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request:", err)
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to post your tweet", resp.StatusCode, http.StatusText(resp.StatusCode), body)
	}
	var tweet map[string]interface{}
	if err := json.Unmarshal(body, &tweet); err != nil {
		return "", fmt.Errorf("error: %v", err)
	}
	tweetID, ok := tweet["data"].(map[string]interface{})["id"].(string)
	if !ok {
		return "", fmt.Errorf("unable to retrieve tweet ID from response: %v", tweet)
	}
	return tweetID, nil
}
func deleteTweet(tweetID string) error {
	apiKey := os.Getenv("TWITTER_API_KEY")
	apiSecretKey := os.Getenv("TWITTER_API_SECRET_KEY")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	// OAuth1 Config setup
	config := oauth1.NewConfig(apiKey, apiSecretKey)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	// HTTP client with OAuth1 credentials
	httpClient := config.Client(oauth1.NoContext, token)
	deleteUrl := fmt.Sprintf(deleteTweetUrl, tweetID)
	// Create the DELETE request to the Twitter API
	req, err := http.NewRequest("DELETE", deleteUrl, nil)
	if err != nil {
		return fmt.Errorf("Error creating request:", err)
	}
	// Send the DELETE request
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending request:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Failed to delete your tweet, statuscode: %d\nResponse: %s\n", resp.StatusCode, body)
	}
	return nil
}
func main() {
	var tweetIDs []string
	tweets := []string{
		"Hii this is the Twitter API!!",
		"Jenish Gaudani here",
		"Hello",
	}
	for _, tweet := range tweets {
		fmt.Println("Posting tweet:", tweet)
		tweetID, err := postTweet(tweet)
		if err != nil {
			fmt.Println("Error with posting your tweet:", err)
			continue
		}
		tweetIDs = append(tweetIDs, tweetID)
		fmt.Println("Tweet posted successfully")
	}
	// Check if there are any tweets posted
	if len(tweetIDs) == 0 {
		fmt.Println("No tweets were posted.")
		return
	}
	tweetIDToDelete := tweetIDs[1]
	fmt.Println("Deleting tweet with ID:", tweetIDToDelete)
	err := deleteTweet(tweetIDToDelete)
	if err != nil {
		fmt.Println("Error deleting tweet:", err)
		return
	}
	fmt.Println("Tweet deleted successfully!")
}
