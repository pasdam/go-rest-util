package restutil

import (
	"encoding/json"
	"net/http"
)

// GetJSON performs a GET request to the specified url. The parsed JSON is
// stored in the value pointed to by v
func GetJSON(url string, headers http.Header, v interface{}) error {
	body, err := Get(url, headers)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
