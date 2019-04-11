package gosalesforcejwt

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestLogIn(t *testing.T) {
	keyBytes, _ := ioutil.ReadFile(os.Getenv("SALESFORCE_KEY_PATH"))
	request, _ := BuildRequest(os.Getenv("SALESFORCE_CLIENT_ID"), os.Getenv("SALESFORCE_USER"), os.Getenv("SALESFORCE_AUDIENCE"))
	signature, _ := SignRequest(keyBytes, request)
	token, err := LogIn(request, signature, os.Getenv("SALESFORCE_ENDPOINT"))
	if err != nil {
		t.Fatal(err)
	}
	log.Println(token)
}
func TestCreateLead(t *testing.T) {
	Init()
	create, err := PostObject("Lead", struct {
		LastName string `json:"LastName"`
		Company  string `json:"Company"`
		Status   string `json:"Status"`
	}{LastName: fmt.Sprintf("Test:%v", time.Now().Unix()), Company: "Test", Status: "0"})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(create.ID)
}
func TestCreateContact(t *testing.T) {

	Init()
	create, err := PostObject("Contact", struct {
		LastName string `json:"LastName"`
		Email    string `json:"Email"`
	}{LastName: "", Email: ""})
	if err != nil {
		log.Println(err, create)
		t.Fatal(err)
	}
	log.Println(create.ID)
}
func TestGetLead(t *testing.T) {
	Init()
	create, err := PostObject("Lead", struct {
		LastName string `json:"LastName"`
		Company  string `json:"Company"`
		Status   string `json:"Status"`
	}{LastName: fmt.Sprintf("Test:%v", time.Now().Unix()), Company: "Test", Status: "0"})
	if err != nil {
		t.Fatal(err)
	}
	lead, err := GetObject("Lead", create.ID, []string{"LastName", "Company", "Status", "IsConverted"})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(lead)
}
func TestPatchLead(t *testing.T) {
	Init()
	create, err := PostObject("Lead", struct {
		LastName string `json:"LastName"`
		Company  string `json:"Company"`
		Status   string `json:"Status"`
	}{LastName: fmt.Sprintf("Test:%v", time.Now().Unix()), Company: "Test", Status: "0"})
	if err != nil {
		t.Fatal(err)
	}
	lead, err := GetObject("Lead", create.ID, []string{"LastName", "Company", "Status", "IsConverted"})
	if err != nil {
		t.Fatal(err)
	}
	updatedName := fmt.Sprintf("%s (updated %v)", lead["LastName"], time.Now().Unix())
	lead["LastName"] = updatedName
	delete(lead, "Id")
	delete(lead, "IsConverted")
	err = PatchObject("Lead", create.ID, lead)
	if err != nil {
		t.Fatal(err)
	}
	lead, err = GetObject("Lead", create.ID, []string{"LastName", "Company", "Status", "IsConverted"})
	if err != nil {
		t.Fatal(err)
	}
	if lead["LastName"].(string) != updatedName {
		t.Fatal("Update failed")
	}
}
func TestSearchLead(t *testing.T) {

	Init()
	name := fmt.Sprintf("Test:%v", time.Now().Unix())
	_, err := PostObject("Lead", struct {
		LastName string `json:"LastName"`
		Company  string `json:"Company"`
		Status   string `json:"Status"`
	}{LastName: name, Company: "Test", Status: "0"})
	if err != nil {
		t.Fatal(err)
	}
	results, err := SearchObject("Lead", name, []string{"Id", "LastName", "Company", "Status", "IsConverted"}, 10)
	if err != nil || len(results) != 1 {
		t.Fatal(err)
	}
	log.Println(results)
}

func TestSearchLeadQuery(t *testing.T) {

	Init()
	name := fmt.Sprintf("Test:%v", time.Now().Unix())
	_, err := PostObject("Lead", struct {
		LastName string `json:"LastName"`
		Company  string `json:"Company"`
		Status   string `json:"Status"`
	}{LastName: name, Company: "Test", Status: "0"})
	if err != nil {
		t.Fatal(err)
	}
	results, err := Query(fmt.Sprintf("select Id,Name,Company from Lead where LastName='%s'", name))
	if err != nil || len(results) != 1 {
		t.Fatal(err)
	}
	log.Println(results)
}

func TestDeleteLead(t *testing.T) {
	Init()
	create, err := PostObject("Lead", struct {
		LastName string `json:"LastName"`
		Company  string `json:"Company"`
		Status   string `json:"Status"`
	}{LastName: fmt.Sprintf("Test:%v", time.Now().Unix()), Company: "Test", Status: "0"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = DeleteObject("Lead", create.ID)
	if err != nil {
		t.Fatal(err)
	}

}
