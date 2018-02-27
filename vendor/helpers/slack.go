package helpers

import (
	"fmt"
	"net/http"
	"strings"
)

// SendCheckout sends a cert request to a slack channel.
func SendCheckout(aid string) error {
	body := strings.NewReader(fmt.Sprintf("{\"text\":\"%s has requested for checkout\"}", aid))
	req, err := http.NewRequest("POST", "https://hooks.slack.com/services/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", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	return nil
}
