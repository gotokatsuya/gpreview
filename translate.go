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

type AccessTokenMessage struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (self *AccessTokenMessage) getExpiresIn() int {
	i, err := strconv.Atoi(self.ExpiresIn)
	if err != nil {
		return 0
	}
	return i
}

type AccessTokenMessageCache struct {
	accessTokenMessage AccessTokenMessage
	updateTime         time.Time
}

func (self *AccessTokenMessageCache) isRequireRenew() bool {
	duration := int(time.Since(self.updateTime)) * int(time.Second)
	if duration >= self.accessTokenMessage.getExpiresIn() {
		return true
	}
	return false
}

func (self *AccessTokenMessageCache) setUpdateTime(time time.Time) {
	self.updateTime = time
}

func (self *AccessTokenMessageCache) getUpdateTime() time.Time {
	return self.updateTime
}

func (self *AccessTokenMessageCache) loadNewAccessTokenMessage() {
	resp, err := http.PostForm(msTranslatorAccessTokenURL,
		url.Values{
			"client_id":     {GPReview.MsTranslatorClientID},
			"client_secret": {GPReview.MsTranslatorClientSecret},
			"scope":         {msTranslatorScope},
			"grant_type":    {msTranslatorGrantType}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&self.accessTokenMessage)
	if err != nil {
		panic(err)
	}
}

func (self *AccessTokenMessageCache) getAccessTokenMessage() AccessTokenMessage {
	if self.isRequireRenew() {
		self.loadNewAccessTokenMessage()
	}
	return self.accessTokenMessage
}

func (self *AccessTokenMessageCache) getAccessToken() string {
	message := self.getAccessTokenMessage()
	return message.AccessToken
}

type Result struct {
	String string `xml:"string"`
}

func Translate(word, from, to string, atmc *AccessTokenMessageCache) string {
	values := url.Values{"text": {word}, "from": {from}, "to": {to}}
	query := values.Encode()
	accessToken := atmc.getAccessToken()

	client := &http.Client{}

	req, err := http.NewRequest("GET", msTranslatorURL+"?"+query, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	bodyString := "<result>" + string(body) + "</result>"

	var res *Result
	err = xml.Unmarshal([]byte(bodyString), &res)
	if err != nil {
		panic(err)
	}
	return res.String
}
