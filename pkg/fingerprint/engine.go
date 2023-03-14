package fingerprint

import (
	"github.com/Escape-Technologies/goctopus/internal/http"
	"github.com/Escape-Technologies/goctopus/pkg/engine"
	log "github.com/sirupsen/logrus"
)

func (fp *fingerprinter) Engine() (string, error) {

	for _, ngin := range engine.Engines {
		for _, imprint := range ngin.Imprints {
			log.Debugf("Trying to match %s with %s", imprint.Query, ngin.Name)
			requestBody := http.GraphqlBodyPayload(imprint.Query)
			resp, err := fp.Client.Post(fp.url, []byte(requestBody))
			log.Debugf("Response: %s", resp.Body)
			if err != nil {
				return "", err
			}
			if imprint.Matcher(resp.Body) {
				return ngin.Name, nil
			}
		}
	}

	return "unknown", nil
}
