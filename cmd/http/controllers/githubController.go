package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GraphQL request payload
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// Expected structure of the GraphQL response
type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []interface{}          `json:"errors"`
}

func GetCommitHistory(c *gin.Context) {
	// GraphQL query
	query := `
		query ($userName: String!) {
			user(userName: $userName) {
				contributionsCollection {
					contributionYears
					contributionCalendar {
						totalContributions
						weeks {
							contributionDays {
								weekday
								date
								contributionCount
								color
							}
						}
					}
				}
			}
		}
	`
	
	userName := c.Query("user")
	if userName == "" {
		log.Error("User parameter is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User parameter is required"})
		return
	}

	// Create the request payload
	requestPayload := GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"userName": userName,
		},
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		log.Error("Error marshaling request payload: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request payload"})
		return
	}

	// Make HTTP POST request to the GitHub GraphQL API
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Error("Error creating HTTP request: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set the authorization header using GitHub token
	githubToken := os.Getenv("GITHUB_ACCESS_TOKEN")
	if githubToken == "" {
		log.Error("GitHub access token is not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GitHub access token is not set"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Error making HTTP request: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to make request"})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Parse JSON response
	var graphqlResponse GraphQLResponse
	err = json.Unmarshal(body, &graphqlResponse)
	if err != nil {
		log.Error("Error unmarshaling response: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Handle potential errors returned by the GraphQL API
	if len(graphqlResponse.Errors) > 0 {
		log.Error("GraphQL API returned errors: ", graphqlResponse.Errors)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch commit history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": graphqlResponse.Data})
}