package gosalesforcejwt

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestLogIn(t *testing.T) {
	godotenv.Load()
	keyBytes, _ := ioutil.ReadFile(os.Getenv("SALESFORCE_KEY_PATH"))

	request, _ := BuildRequest(os.Getenv("SALESFORCE_CLIENT_ID"), os.Getenv("SALESFORCE_USER"), os.Getenv("SALESFORCE_AUDIENCE"))

	signature, _ := SignRequest(keyBytes, request)

	token, err, errr := LogIn(request, signature, os.Getenv("SALESFORCE_ENDPOINT"))

	log.Println(err, errr, token)

}
