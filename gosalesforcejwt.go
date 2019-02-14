package gosalesforcejwt

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

//SalesforceTokenResponse represents the standars succesful authentication response
type SalesforceTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	InstanceURL string `json:" instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
}

//SalesforceErrorResponse represents a standard flow error response
type SalesforceErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

//JWTHeader represents Salesforce required header fields in a standard JWT header format
type JWTHeader struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

//JWTClaims repesents Salesforce required claims in standard JWT claim format
type JWTClaims struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp,number"`
}

//JWTFlowLogIn returns either a SalesforceTokenResponse or an error and possibly a SalesforceErrorResponse to help with debugging if communication was succesful but configuration issues were found
func JWTFlowLogIn(sandbox bool) (*SalesforceTokenResponse, error, *SalesforceErrorResponse) {

	clientID := os.Getenv("CLIENT_ID")
	user := os.Getenv("SF_USER")
	var aud string

	if sandbox {
		aud = "https://test.salesforce.com"

	} else {
		aud = "https://login.salesforce.com"

	}

	exp := time.Now().Add(time.Minute * time.Duration(5)).Unix()

	header := JWTHeader{Typ: "JWT", Alg: "RS256"}

	claims := JWTClaims{Iss: clientID, Sub: user, Aud: aud, Exp: exp}

	jh, err := json.Marshal(header)

	if err != nil {
		return nil, err, nil
	}
	jc, err := json.Marshal(claims)

	if err != nil {
		return nil, err, nil
	}

	request := encodeBase64URL(jh) + "." + encodeBase64URL(jc)

	s, err := sign(request)

	signature := encodeBase64URL(s)

	signedRequest := request + "." + signature

	body := "grant_type=" + url.QueryEscape("urn:ietf:params:oauth:grant-type:jwt-bearer") + "&assertion=" + signedRequest

	authRequest, err := http.NewRequest("POST", os.Getenv("INSTANCE")+"/services/oauth2/token", bytes.NewBuffer([]byte(body)))

	if err != nil {
		return nil, err, nil
	}
	authRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	authResponse, err := http.DefaultClient.Do(authRequest)

	if err != nil {
		return nil, err, nil
	}

	if authResponse.StatusCode != 200 {

		errorResponse := &SalesforceErrorResponse{}

		json.NewDecoder(authResponse.Body).Decode(errorResponse)

		return nil, errors.New(errorResponse.ErrorDescription), errorResponse

	}

	token := &SalesforceTokenResponse{}

	json.NewDecoder(authResponse.Body).Decode(token)

	return token, nil, nil

}

func sign(request string) ([]byte, error) {

	keyBytes, err := ioutil.ReadFile(os.Getenv("KEY_PATH"))

	if err != nil {
		return nil, err
	}

	block, rest := pem.Decode([]byte(keyBytes))

	if len(rest) != 0 {
		return nil, errors.New("Invalid key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	sum := sha256.Sum256([]byte(request))

	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sum[:])

	if err != nil {
		return nil, err
	}
	return sigBytes, nil

}

func encodeBase64URL(data []byte) string {

	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")

}
