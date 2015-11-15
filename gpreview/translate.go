package gpreview

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	msTranslatorScope          = "http://api.microsofttranslator.com"
	msTranslatorAccessTokenURL = "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"
	msTranslatorURL            = "http://api.microsofttranslator.com/V2/Http.svc/Translate"
	msTranslatorGrantType      = "client_credentials"
)

type MsAccessTokenMessage struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (m *MsAccessTokenMessage) getExpiresIn() int {
	i, err := strconv.Atoi(m.ExpiresIn)
	if err != nil {
		return 0
	}
	return i
}

type MsAccessTokenMessageCache struct {
	accessTokenMessage MsAccessTokenMessage
	updateTime         time.Time
}

func (ms *MsAccessTokenMessageCache) isRequireRenew() bool {
	if len(ms.accessTokenMessage.AccessToken) <= 0 {
		return true
	}
	duration := int(time.Since(ms.updateTime)) * int(time.Second)
	if duration >= ms.accessTokenMessage.getExpiresIn() {
		return true
	}
	return false
}

func (ms *MsAccessTokenMessageCache) loadNewAccessTokenMessage() error {
	resp, err := http.PostForm(msTranslatorAccessTokenURL,
		url.Values{
			"client_id":     {GPReview.MsTranslatorClientID},
			"client_secret": {GPReview.MsTranslatorClientSecret},
			"scope":         {msTranslatorScope},
			"grant_type":    {msTranslatorGrantType}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&ms.accessTokenMessage)
	if err != nil {
		return err
	}
	ms.updateTime = time.Now()
	return nil
}

func (ms *MsAccessTokenMessageCache) getAccessTokenMessage() MsAccessTokenMessage {
	if ms.isRequireRenew() {
		ms.loadNewAccessTokenMessage()
	}
	return ms.accessTokenMessage
}

func (ms *MsAccessTokenMessageCache) getAccessToken() string {
	message := ms.getAccessTokenMessage()
	return message.AccessToken
}

type Result struct {
	String string `xml:"string"`
}

func Translate(word, from, to string, atmc *MsAccessTokenMessageCache) (string, error) {
	values := url.Values{"appId": {""}, "text": {word}, "from": {from}, "to": {to}}
	query := values.Encode()
	accessToken := atmc.getAccessToken()

	client := &http.Client{}

	req, err := http.NewRequest("GET", msTranslatorURL+"?"+query, nil)
	if err != nil {
		return word, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return word, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	bodyString := "<result>" + string(body) + "</result>"

	var res *Result
	err = xml.Unmarshal([]byte(bodyString), &res)
	if err != nil {
		return word, err
	}
	return res.String, nil
}
