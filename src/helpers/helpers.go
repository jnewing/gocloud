package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// FetchExtIP
// So this does just that (^ see above) it uses an external API to get your current IP (or the IP of the calling machiene)
// it honestly really doesn't matter what you use here as long as it returns the IP you want in the form of a string.
// So if you want to modify this go for it.
func FetchExtIP() (*string, error) {
	resp, err := http.Get("https://api.bigdatacloud.net/data/client-ip")

	if err != nil {
		return nil, fmt.Errorf("error while fetching external ip")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error while reading response")
	}

	var result map[string]string
	json.Unmarshal([]byte(body), &result)

	ipaddress := string(result["ipString"])

	return &ipaddress, nil
}
