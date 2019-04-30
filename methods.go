package gosalesforcejwt

import (
	"bytes"
	"encoding/json"
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

type SalesforceAPIResponse struct {
	ID      string   `json:"id"`
	Errors  []string `json:"errors"`
	Success bool     `json:"success"`
}

func Init() error {
	keyBytes, err := ioutil.ReadFile(os.Getenv("SALESFORCE_KEY_PATH"))
	if err != nil {
		log.Panicln(err)
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
	reqURL := fmt.Sprintf("%s/services/data/v45.0/sobjects/%s", os.Getenv("SALESFORCE_ENDPOINT"), objectName)
	jc, _ := json.Marshal(data)
	r, _ := http.NewRequest("POST", reqURL, bytes.NewBuffer(jc))
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 400 {
		bs, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("%s:Status %v:%s", objectName, resp.StatusCode, string(bs))
	}
	respBody := &SalesforceAPIResponse{}
	json.NewDecoder(resp.Body).Decode(respBody)

	return respBody, nil

}
func GetObject(objectName string, ID string, fields []string) (map[string]interface{}, error) {
	reqURL := fmt.Sprintf("%s/services/data/v45.0/sobjects/%s/%s/?fields=%s", os.Getenv("SALESFORCE_ENDPOINT"), objectName, ID, strings.Join(fields, ","))
	r, _ := http.NewRequest("GET", reqURL, nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Accept", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 404 || resp.StatusCode == 400 {
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s:Status %v:%s", objectName, resp.StatusCode, string(bs))
	}
	respBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&respBody)
	return respBody, nil
}

func PatchObject(objectName string, ID string, data interface{}) error {
	reqURL := fmt.Sprintf("%s/services/data/v45.0/sobjects/%s/%s", os.Getenv("SALESFORCE_ENDPOINT"), objectName, ID)
	jc, _ := json.Marshal(data)
	r, _ := http.NewRequest("PATCH", reqURL, bytes.NewBuffer(jc))
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 || resp.StatusCode == 400 {
		bs, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%s:Status %v:%s", objectName, resp.StatusCode, string(bs))
	}
	return nil
}

func SearchObject(objectName string, query string, fields []string, limit int) ([]interface{}, error) {
	reqURL := fmt.Sprintf("%s/services/data/v45.0/parameterizedSearch/?q=%s&sobject=%s&%s.fields=%s&%s.limit=%v", os.Getenv("SALESFORCE_ENDPOINT"), url.QueryEscape(query), objectName, objectName, url.QueryEscape(strings.Join(fields, ",")), objectName, limit)
	r, _ := http.NewRequest("GET", reqURL, nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Accept", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bs, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("%s:Status %v: %s with query %s", objectName, resp.StatusCode, string(bs), reqURL)
	}
	responseMap := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseMap)
	if responseMap["searchRecords"] == nil {
		return nil, fmt.Errorf("Unexpected response body %v", responseMap)
	}
	return responseMap["searchRecords"].([]interface{}), nil
}

func Query(query string) ([]interface{}, error) {
	reqURL := fmt.Sprintf("%s/services/data/v45.0/query/?q=%s", os.Getenv("SALESFORCE_ENDPOINT"), url.QueryEscape(query))
	r, _ := http.NewRequest("GET", reqURL, nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Accept", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		bs, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("Status %v: %s with query %s", resp.StatusCode, string(bs), reqURL)
	}
	responseMap := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&responseMap)
	if responseMap["records"] == nil {
		return nil, fmt.Errorf("Unexpected response body %v", responseMap)
	}
	return responseMap["records"].([]interface{}), nil
}
func DeleteObject(objectName string, ID string) (map[string]interface{}, error) {
	reqURL := fmt.Sprintf("%s/services/data/v45.0/sobjects/%s/%s", os.Getenv("SALESFORCE_ENDPOINT"), objectName, ID)
	r, _ := http.NewRequest("DELETE", reqURL, nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", salesforceToken.AccessToken))
	r.Header.Set("Accept", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 204 {
		bs, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%s:Status %v:%s", objectName, resp.StatusCode, string(bs))
	}
	respBody := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&respBody)
	return respBody, nil
}
