package gosalesforcejwt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var salesforceToken *SalesforceTokenResponse
var client *http.Client
var logging bool

type SalesforceAPIResponse struct {
	ID      string   `json:"id"`
	Errors  []string `json:"errors"`
	Success bool     `json:"success"`
}

func Init(logRequests bool) error {
	logging = logRequests
	keyBytes, err := ioutil.ReadFile(os.Getenv("SALESFORCE_KEY_PATH"))
	if err != nil {
		return err
	}
	request, err := BuildRequest(os.Getenv("SALESFORCE_CLIENT_ID"), os.Getenv("SALESFORCE_USER"), os.Getenv("SALESFORCE_AUDIENCE"))
	signature, err := SignRequest(keyBytes, request)
	token, err := LogIn(request, signature, os.Getenv("SALESFORCE_ENDPOINT"))
	if err != nil {
		return err
	}
	salesforceToken = token
	client = &http.Client{}
	return nil
}

func PostObject(objectName string, data interface{}) (*SalesforceAPIResponse, error) {
	var bs []byte
	var loggable string

	reqURL := fmt.Sprintf("%s/services/data/v45.0/sobjects/%s", os.Getenv("SALESFORCE_ENDPOINT"), objectName)
	jc, _ := json.Marshal(data)
	r, _ := http.NewRequest("POST", reqURL, bytes.NewBuffer(jc))
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if logging {
		bs, loggable, err = logResponse(resp, reqURL)
		log.Println(loggable)
	} else {
		bs, err = ioutil.ReadAll(resp.Body)
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(loggable)
	}
	respBody := &SalesforceAPIResponse{}
	json.NewDecoder(bytes.NewBuffer(bs)).Decode(&respBody)

	return respBody, nil

}
func GetObject(objectName string, ID string, fields []string) (map[string]interface{}, error) {
	var bs []byte
	var loggable string

	reqURL := fmt.Sprintf("%s/services/data/v45.0/sobjects/%s/%s/?fields=%s", os.Getenv("SALESFORCE_ENDPOINT"), objectName, ID, strings.Join(fields, ","))
	r, _ := http.NewRequest("GET", reqURL, nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Accept", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if logging {
		bs, loggable, err = logResponse(resp, reqURL)
		log.Println(loggable)
	} else {
		bs, err = ioutil.ReadAll(resp.Body)
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(loggable)
	}
	respBody := make(map[string]interface{})
	json.NewDecoder(bytes.NewBuffer(bs)).Decode(&respBody)
	return respBody, nil
}

func PatchObject(objectName string, ID string, data interface{}) error {
	var loggable string

	reqURL := fmt.Sprintf("%s/services/data/v45.0/sobjects/%s/%s", os.Getenv("SALESFORCE_ENDPOINT"), objectName, ID)
	jc, _ := json.Marshal(data)
	r, _ := http.NewRequest("PATCH", reqURL, bytes.NewBuffer(jc))
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return err
	}

	if logging {
		_, loggable, err = logResponse(resp, reqURL)
	}

	if resp.StatusCode > 299 {
		return errors.New(loggable)
	}
	return nil
}

func SearchObject(objectName string, query string, fields []string, limit int) ([]interface{}, error) {
	var bs []byte
	var loggable string

	reqURL := fmt.Sprintf("%s/services/data/v45.0/parameterizedSearch/?q=%s&sobject=%s&%s.fields=%s&%s.limit=%v", os.Getenv("SALESFORCE_ENDPOINT"), url.QueryEscape(query), objectName, objectName, url.QueryEscape(strings.Join(fields, ",")), objectName, limit)
	r, _ := http.NewRequest("GET", reqURL, nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Accept", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if logging {
		bs, loggable, err = logResponse(resp, reqURL)
		log.Println(loggable)
	} else {
		bs, err = ioutil.ReadAll(resp.Body)
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(loggable)
	}
	respBody := make(map[string]interface{})
	json.NewDecoder(bytes.NewBuffer(bs)).Decode(&respBody)

	if sr, e := respBody["searchRecords"]; !e || sr == nil {
		return nil, errors.New(loggable)
	} else if sa, is := sr.([]interface{}); is {
		return sa, nil

	} else {
		return nil, errors.New(loggable)
	}

}

func Query(query string) ([]interface{}, error) {
	var bs []byte
	var loggable string

	reqURL := fmt.Sprintf("%s/services/data/v45.0/query/?q=%s", os.Getenv("SALESFORCE_ENDPOINT"), url.QueryEscape(query))
	r, _ := http.NewRequest("GET", reqURL, nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Accept", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if logging {
		bs, loggable, err = logResponse(resp, reqURL)
		log.Println(loggable)
	} else {
		bs, err = ioutil.ReadAll(resp.Body)
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(loggable)
	}
	respBody := make(map[string]interface{})
	json.NewDecoder(bytes.NewBuffer(bs)).Decode(&respBody)

	if sr, e := respBody["records"]; !e || sr == nil {
		return nil, errors.New(loggable)
	} else if sa, is := sr.([]interface{}); is {
		return sa, nil

	} else {
		return nil, errors.New(loggable)
	}
}
func DeleteObject(objectName string, ID string) (map[string]interface{}, error) {
	var bs []byte
	var loggable string
	reqURL := fmt.Sprintf("%s/services/data/v45.0/sobjects/%s/%s", os.Getenv("SALESFORCE_ENDPOINT"), objectName, ID)
	r, _ := http.NewRequest("DELETE", reqURL, nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Accept", "application/json")

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if logging {
		bs, loggable, err = logResponse(resp, reqURL)
		log.Println(loggable)
	} else {
		bs, err = ioutil.ReadAll(resp.Body)
	}

	if resp.StatusCode > 299 {
		return nil, errors.New(loggable)
	}
	respBody := make(map[string]interface{})
	json.NewDecoder(bytes.NewBuffer(bs)).Decode(&respBody)
	return respBody, nil
}

func logResponse(res *http.Response, url string) ([]byte, string, error) {
	var body interface{}
	var cp []byte

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	if len(bs) >= 0 {
		cp = make([]byte, len(bs))
		copy(cp, bs)

		err = json.NewDecoder(bytes.NewBuffer(bs)).Decode(&body)
		if err != nil {
			return nil, "", err
		}
	}

	js, err := json.Marshal(map[string]interface{}{"url": url, "status_code": res.StatusCode, "body": body})
	if err != nil {
		return nil, "", err
	}

	return cp, string(js), nil

}
