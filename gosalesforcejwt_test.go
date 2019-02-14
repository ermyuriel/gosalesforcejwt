package gosalesforcejwt

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestLogIn(t *testing.T) {
	godotenv.Load()
	token, err, errorResponse := JWTFlowLogIn(false)

	if err != nil {
		t.Fail()
	}

	if errorResponse != nil {
		t.Fail()
	}

	j, _ := json.MarshalIndent(token, " ", "")

	log.Println(string(j))

}
