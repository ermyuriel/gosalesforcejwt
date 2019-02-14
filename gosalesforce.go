package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type JWTHeader struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
}

type JWTClaims struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp,number"`
}

func JWTFlowLogIn(sandbox bool) (string, error) {

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
		return "", err
	}
	jc, err := json.Marshal(claims)

	if err != nil {
		return "", err
	}

	request := EncodeBase64URL(jh) + "." + EncodeBase64URL(jc)

	s := sign(request)

	signature := EncodeBase64URL(s)

	signedRequest := request + "." + signature

	log.Println(signedRequest)

	body := "grant_type=" + url.QueryEscape("urn:ietf:params:oauth:grant-type:jwt-bearer") + "&assertion=" + signedRequest

	authRequest, err := http.NewRequest("POST", os.Getenv("INSTANCE")+"/services/oauth2/token", bytes.NewBuffer([]byte(body)))

	if err != nil {
		log.Panicln(err)
	}
	authRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	r, _ := httputil.DumpRequest(authRequest, true)
	log.Println(string(r))

	authResponse, err := http.DefaultClient.Do(authRequest)

	if err != nil {
		log.Panicln(err)
	}

	r, _ = httputil.DumpResponse(authResponse, true)
	log.Println(string(r))

	return "", nil

}

func main() {

	godotenv.Load()
	JWTFlowLogIn(false)

}

func sign(request string) []byte {

	keyBytes, err := ioutil.ReadFile(os.Getenv("KEY_PATH"))
	if err != nil {
		log.Panicln(err)
	}
	block, rest := pem.Decode([]byte(keyBytes))
	if len(rest) != 0 {
		log.Panicln("Invalid key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Panicln(err)
	}

	sum := sha256.Sum256([]byte(request))

	if sigBytes, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sum[:]); err == nil {
		return sigBytes
	} else {
		log.Panicln(err)
	}

	return []byte{}
}
