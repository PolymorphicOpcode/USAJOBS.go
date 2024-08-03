package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

// Define the job structure based on the API response
type Job struct {
	JobTitle  string `json:"PositionTitle"`
	Agency    string `json:"OrganizationName"`
	OpenDate  string `json:"PublicationStartDate"`
	Link      string `json:"ApplyURI"`
	ControlID string `json:"MatchedObjectId"`
}

type APIResponse struct {
	SearchResult struct {
		SearchResultItems []struct {
			MatchedObjectDescriptor Job `json:"MatchedObjectDescriptor"`
		} `json:"SearchResultItems"`
	} `json:"SearchResult"`
}

var ctx = context.Background()

func main() {
	// Retrieve email and API key from environment variables
	email := os.Getenv("USAJOBS_EMAIL")
	apiKey := os.Getenv("USAJOBS_API_KEY")

	if email == "" || apiKey == "" {
		log.Fatal("Environment variables USAJOBS_EMAIL and USAJOBS_API_KEY must be set")
	}

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Adjust this if your Redis server is running elsewhere
		DB:   0,                // Use default DB
	})

	// Set up the API request
	apiURL := "https://data.usajobs.gov/api/search"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set the required headers
	req.Header.Set("User-Agent", email)
	req.Header.Set("Authorization-Key", apiKey)

	// Add query parameters
	q := req.URL.Query()
	q.Add("Keyword", "CYBER")
	req.URL.RawQuery = q.Encode()

	// Perform the API request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to perform API request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the JSON response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Process and store job listings
	for _, item := range apiResp.SearchResult.SearchResultItems {
		job := item.MatchedObjectDescriptor
		if err := rdb.Get(ctx, job.ControlID).Err(); err == redis.Nil {
			jobInfo := fmt.Sprintf("Title: %s, Agency: %s, Open Date: %s, Link: %s", job.JobTitle, job.Agency, job.OpenDate, job.Link)
			rdb.Set(ctx, job.ControlID, jobInfo, 0)
			fmt.Println(jobInfo)
		}
	}
}
