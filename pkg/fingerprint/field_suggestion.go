// This is WIP, not implemented yet
package fingerprint

import (
	"encoding/json"
	"regexp"

	"github.com/Escape-Technologies/goctopus/internal/http"
	"github.com/Escape-Technologies/goctopus/internal/utils"
	log "github.com/sirupsen/logrus"
)

// @todo this also probably needs to match if the query was right
var (
	SuggestionRegexp *regexp.Regexp
)

func init() {
	SuggestionRegexp = regexp.MustCompile(`.*Did you mean.*`)
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

func (fp *fingerprinter) FieldSuggestion() (bool, error) {
	for _, word := range *utils.Wordlist {
		body := makePayload(word)
		res, err := fp.Client.Post(fp.url, body)
		log.Debugf("Response from %v: %v", fp.url, res.StatusCode)
		if err != nil {
			log.Debugf("Error from %v: %v", fp.url, err)
			return false, err
		}
		if IsSuggestionResponse(res) {
			return true, nil
		}
	}
	return false, nil
}
