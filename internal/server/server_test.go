package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/StStep/go-test-server/internal/auth"
)

func TestGetHome(t *testing.T) {
	// Setup
	svr := Start()

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

	resp, err := client.Do(r)
	if err != nil {
		svr.Close()
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is %v not %v", resp.Status, http.StatusOK)
	}

	// Decode Body
	var st Status
	if json.NewDecoder(resp.Body).Decode(&st) != nil {
		svr.Close()
		panic(err)
	}

	// Check Msg Body
	if st.Status != "Ready" {
		t.Errorf("Server Status is %v not %v", st.Status, "Ready")
	}
	if st.PeopleNumber != 3 {
		t.Errorf("PeopleNumber is %v not %v", st.PeopleNumber, 3)
	}

	// Tear-down
	svr.Close()
}

func TestPostLogin(t *testing.T) {
	// Setup
	svr := Start()

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	r_people, _ := http.NewRequest("GET", "http://localhost:8080/people", nil)
	p_login, _ := http.NewRequest("POST", "http://localhost:8080/login", nil)

	// requesting people should fail before logging in
	resp, err := client.Do(r_people)
	if err != nil {
		svr.Close()
		panic(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("%v is allowed to be viewed without authorization", r_people.URL)
	}

	// Posting login should return valid token
	resp, err = client.Do(p_login)
	if err != nil {
		svr.Close()
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is %v not %v", resp.Status, http.StatusOK)
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		svr.Close()
		panic(err)
	}
	token := string(respData)
	err = auth.VerifyToken(token)
	if err != nil {
		t.Errorf("Received token is invalid with error: %v", err)
	}

	// Should succeed with new token
	r_people.SetBasicAuth("token", token)
	resp, err = client.Do(r_people)
	if err != nil {
		svr.Close()
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("%v failed authorization with valid token", r_people.URL)
	}

	// Tear-down
	svr.Close()
}
