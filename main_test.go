package main

import (
        "net/http"
        "testing"
        "time"
	"encoding/json"
)

func TestGetHome(t *testing.T) {
        go startServer()
        client := &http.Client{
                Timeout: 1 * time.Second,
        }

        r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

        resp, err := client.Do(r)
        if err != nil {
                panic(err)
        }
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code is %v not %v", resp.Status, http.StatusOK)
	}

	// Decode Body
	var st Status
	if json.NewDecoder(resp.Body).Decode(&st) != nil {
                panic(err)
        }

	// Check Msg Body
	if st.Status != "Ready" {
		t.Errorf("Server Status is %v not %v", st.Status, "Ready")
	}
	if st.PeopleNumber != 3 {
		t.Errorf("PeopleNumber is %v not %v", st.PeopleNumber, 3)
	}
}
