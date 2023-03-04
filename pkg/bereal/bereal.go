package bereal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ConfigAPIBaseURL       = "https://mobile.bereal.com/api"
	ConfigGoogleAPIBaseURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty"
	ConfigGoogleAPIKey     = "AIzaSyDwjfEeparokD7sXPVQli9NsTuhT6fJ6iA" //nolint
	ConfigRequestHeaders   = map[string]string{
		"x-firebase-client":          "apple-platform/ios apple-sdk/19F64 appstore/true deploy/cocoapods device/iPhone9,1 fire-abt/8.15.0 fire-analytics/8.15.0 fire-auth/8.15.0 fire-db/8.15.0 fire-dl/8.15.0 fire-fcm/8.15.0 fire-fiam/8.15.0 fire-fst/8.15.0 fire-fun/8.15.0 fire-install/8.15.0 fire-ios/8.15.0 fire-perf/8.15.0 fire-rc/8.15.0 fire-str/8.15.0 firebase-crashlytics/8.15.0 os-version/14.7.1 xcode/13F100",
		"user-agent":                 "FirebaseAuth.iOS/8.15.0 AlexisBarreyat.BeReal/0.22.4 iPhone/14.7.1 hw/iPhone9_1",
		"x-ios-bundle-identifier":    "AlexisBarreyat.BeReal",
		"x-firebase-client-log-type": "0",
		"x-client-version":           "iOS/FirebaseSDK/8.15.0/FirebaseCore-iOS",
	}
	ConfigIOSString = "AEFDNu9QZBdycrEZ8bM_2-Ei5kn6XNrxHplCLx2HYOoJAWx-uSYzMldf66-gI1vOzqxfuT4uJeMXdreGJP5V1pNen_IKJVED3EdKl0ldUyYJflW5rDVjaQiXpN0Zu2BNc1c"
	ConfigIOSSecret = "KKwuB8YqwuM3ku0z" //nolint
)

type BeReal struct {
	sessionInfo  string
	idToken      string
	refreshToken string
	localID      string

	Debug bool
}

type erroResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
		} `json:"errors"`
	} `json:"error"`
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

	if resp.StatusCode != 200 {
		var errorResp erroResponse
		err = json.Unmarshal(bodyText, &errorResp)
		if err != nil {
			return err
		}
		return errors.New(errorResp.Error.Message)
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
	IOSSecret   string `json:"iosSecret"`
}

func (b *BeReal) SendAuthMessage(phoneNumber string) error {
	var response sendAuthMessageResponse
	err := request(
		fmt.Sprintf("%s/sendVerificationCode?key=%s", ConfigGoogleAPIBaseURL, ConfigGoogleAPIKey),
		"POST",
		sendAuthMessageRequest{
			PhoneNumber: phoneNumber,
			IOSReceipt:  ConfigIOSString,
			IOSSecret:   ConfigIOSSecret,
		},
		ConfigRequestHeaders,
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
		fmt.Sprintf("%s/verifyPhoneNumber?key=%s", ConfigGoogleAPIBaseURL, ConfigGoogleAPIKey),
		"POST",
		verifyAuthMessageRequest{
			SessionInfo: b.sessionInfo,
			Code:        code,
			Operation:   "SIGN_UP_OR_IN",
		},
		ConfigRequestHeaders,
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
		fmt.Sprintf("%s/feeds/memories", ConfigAPIBaseURL),
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
