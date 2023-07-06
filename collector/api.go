package collector

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//APICall is what we use to query the HPMSA API
func APICall(client *http.Client, hostName string, sessionKey string, path string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/show/%s", hostName, path), nil)
	req.Header.Add("sessionKey", sessionKey)
	if err != nil {
		log.Print("ERROR: ", err)
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("ERROR: ", err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("ERROR: ", err)
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}
