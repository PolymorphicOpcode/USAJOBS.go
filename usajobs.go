package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Define the job structure based on the API response
type Job struct {
	//PositionID           string `json:"PositionID"`
	JobTitle             string `json:"PositionTitle"`
	Agency               string `json:"OrganizationName"`
	PublicationStartDate string `json:"PublicationStartDate"`
	ApplicationCloseDate string `json:"ApplicationCloseDate"`
	PositionURI          string `json:"PositionURI"`
	ControlID            string `json:"MatchedObjectId"`
}

type APIResponse struct {
	SearchResult struct {
		SearchResultItems []struct {
			MatchedObjectDescriptor Job `json:"MatchedObjectDescriptor"`
		} `json:"SearchResultItems"`
	} `json:"SearchResult"`
}

func main() {
	// Retrieve credentials from environment variables or prompt the user
	email, apiKey := getCredentials()

	// Set up the API request
	apiURL := "https://data.usajobs.gov/api/search"

	// Create the HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	// Set the required headers
	req.Header.Set("User-Agent", email)
	req.Header.Set("Authorization-Key", apiKey)

	// Add query parameters
	q := req.URL.Query()
	q.Add("HiringPath", "student")
	q.Add("SortField", "opendate")
	q.Add("SortDirection", "desc")
	q.Add("ResultsPerPage", "500")
	// Remove for all results instead of this week
	q.Add("DatePosted", "7")

	req.URL.RawQuery = q.Encode()

	// Perform the API request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to perform API request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Prepare and display job listings in a table
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Open Date", "Title", "Agency", "Close Date", "Link"})

	for _, item := range apiResp.SearchResult.SearchResultItems {
		job := item.MatchedObjectDescriptor
		t.AppendRow(table.Row{
			//job.PositionID,
			job.PublicationStartDate,
			job.JobTitle,
			job.Agency,
			job.ApplicationCloseDate,
			job.PositionURI,
		})
	}
	t.Render()
}

// getCredentials retrieves credentials from environment variables or prompts the user
func getCredentials() (string, string) {
	email := os.Getenv("USAJOBS_EMAIL")
	apiKey := os.Getenv("USAJOBS_API_KEY")

	if email == "" || apiKey == "" {
		reader := bufio.NewReader(os.Stdin)

		if email == "" {
			fmt.Print("Enter your USAJobs Email: ")
			email, _ = reader.ReadString('\n')
			email = strings.TrimSpace(email)
		}

		if apiKey == "" {
			fmt.Print("Enter your USAJobs API Key: ")
			apiKey, _ = reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKey)
		}
	}

	return email, apiKey
}
