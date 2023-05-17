package suggestion

import (
	"encoding/json"
	"regexp"

	"github.com/Escape-Technologies/goctopus/internal/utils"
	"github.com/Escape-Technologies/goctopus/pkg/http"
	log "github.com/sirupsen/logrus"
)

var (
	SuggestionRegexp *regexp.Regexp
)

func init() {
	SuggestionRegexp = regexp.MustCompile(`(?i).*did you mean.*`)
}

func MatchFieldSuggestionRegex(message string) bool {
	return SuggestionRegexp.MatchString(message)
}

func IsSuggestionResponse(resp *http.Response) bool {
	body := resp.Body

	type Response struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	var result Response
	if err := json.Unmarshal(*body, &result); err != nil {
		return false
	}
	if len(result.Errors) == 0 {
		return false
	}
	for _, err := range result.Errors {
		if err.Message == "" {
			continue
		}
		if MatchFieldSuggestionRegex(err.Message) {
			return true
		}
	}
	return false
}

func makePayload(word string) []byte {
	return []byte(`{"query": "{` + word + `}"}`)
}

func FingerprintFieldSuggestion(url string, client http.Client) bool {
	for _, word := range *utils.Wordlist {
		body := makePayload(word)
		res, err := client.Post(url, body)
		if err != nil {
			log.Debugf("Error from %v: %v", url, err)
			return false
		}
		log.Debugf("Response from %v: %v", url, res.StatusCode)
		if IsSuggestionResponse(res) {
			return true
		}
	}
	return false
}
