package gpreview

import (
	"encoding/json"
	"net/http"
	nurl "net/url"
	"strings"
)

// SlackData ...
type SlackData struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
}

// PostSlack ...
func PostSlack(d SlackData, url string) error {
	// json
	jsonBody, err := json.Marshal(d)
	if err != nil {
		return err
	}

	// parameters
	v := nurl.Values{}
	v.Add("payload", string(jsonBody))

	// request
	req, err := http.NewRequest("POST", url, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// do request
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
