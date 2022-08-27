package bereal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	CONFIG_BEREAL_API_BASE_URL = "https://mobile.bereal.com/api"
	CONFIG_GOOGLE_API_BASE_URL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty"
	CONFIG_GOOGLE_API_KEY      = "AIzaSyDwjfEeparokD7sXPVQli9NsTuhT6fJ6iA"
	CONFIG_REQUEST_HEADERS     = map[string]string{"user-agent": "AlexisBarreyat.BeReal/0.23.2 iPhone/16.0 hw/iPhone13_2", "x-ios-bundle-identifier": "AlexisBarreyat.BeReal"}
	CONFIG_IOS_STRING          = "AEFDNu9QZBdycrEZ8bM_2-Ei5kn6XNrxHplCLx2HYOoJAWx-uSYzMldf66-gI1vOzqxfuT4uJeMXdreGJP5V1pNen_IKJVED3EdKl0ldUyYJflW5rDVjaQiXpN0Zu2BNc1c"
)

type BeReal struct {
	sessionInfo  string
	idToken      string
	refreshToken string
	localID      string

	Debug bool
}

func request(
	url string,
	method string,
	payload interface{},
	headers map[string]string,
	response interface{},
) error {
	client := &http.Client{}

	var req *http.Request
	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		bodyReader := bytes.NewReader(body)
		req, err = http.NewRequest(method, url, bodyReader)
		if err != nil {
			return err
		}
	} else {
		var err error
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return err
		}
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyText, &response)
	if err != nil {
		return err
	}
	return nil
}

type sendAuthMessageResponse struct {
	SessionInfo string `json:"sessionInfo"`
}

type sendAuthMessageRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	IOSReceipt  string `json:"iosReceipt"`
}

func (b *BeReal) SendAuthMessage(phoneNumber string) error {

	var response sendAuthMessageResponse
	err := request(
		fmt.Sprintf("%s/sendVerificationCode?key=%s", CONFIG_GOOGLE_API_BASE_URL, CONFIG_GOOGLE_API_KEY),
		"POST",
		sendAuthMessageRequest{
			PhoneNumber: phoneNumber,
			IOSReceipt:  CONFIG_IOS_STRING,
		},
		CONFIG_REQUEST_HEADERS,
		&response,
	)
	if err != nil {
		return err
	}
	if b.Debug {
		log.Println("Response Info: " + response.SessionInfo)
	}
	b.sessionInfo = response.SessionInfo
	return nil
}

type verifyAuthMessageResponse struct {
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	LocalID      string `json:"localId"`
}

type verifyAuthMessageRequest struct {
	SessionInfo string `json:"sessionInfo"`
	Code        string `json:"code"`
	Operation   string `json:"operation"`
}

func (b *BeReal) VerifyAuthMessage(code string) error {

	var response verifyAuthMessageResponse
	err := request(
		fmt.Sprintf("%s/verifyPhoneNumber?key=%s", CONFIG_GOOGLE_API_BASE_URL, CONFIG_GOOGLE_API_KEY),
		"POST",
		verifyAuthMessageRequest{
			SessionInfo: b.sessionInfo,
			Code:        code,
			Operation:   "SIGN_UP_OR_IN",
		},
		CONFIG_REQUEST_HEADERS,
		&response,
	)
	if err != nil {
		return err
	}

	if b.Debug {
		log.Println("Response Info: " + response.IDToken)
	}

	b.idToken = response.IDToken
	b.refreshToken = response.RefreshToken
	b.localID = response.LocalID
	return nil
}

type getMemoriesResponse struct {
	Data []Memory `json:"data"`
}

func (b *BeReal) GetMemories() ([]Memory, error) {

	var response getMemoriesResponse
	err := request(
		fmt.Sprintf("%s/feeds/memories", CONFIG_BEREAL_API_BASE_URL),
		"GET",
		nil,
		map[string]string{
			"authorization": b.idToken,
		},

		&response,
	)
	if err != nil {
		return []Memory{}, err
	}

	return response.Data, nil
}
