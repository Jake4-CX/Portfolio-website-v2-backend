package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

func VerifyRecaptcha(recaptchaToken string) bool {

	payload := url.Values{
		"secret":   {os.Getenv("RECAPTCHA_SECRET_KEY")},
		"response": {recaptchaToken},
	}

	// Make POST request
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", payload)
	if err != nil {
		log.Errorf("Error making POST request for recaptcha: %v", err)
		return false
	}

	defer resp.Body.Close() // Close the response body when the function returns

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var response struct {
		Success     bool     `json:"success"`
		ChallengeTs string   `json:"challenge_ts"`
		Hostname    string   `json:"hostname"`
		Score       float64  `json:"score"`
		Action      string   `json:"action"`
		ErrorCodes  []string `json:"error-codes"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false
	}

	return response.Success
}
