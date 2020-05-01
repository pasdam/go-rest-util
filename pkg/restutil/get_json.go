package restutil

import "encoding/json"

// GetJSON performs a GET request to the specified url. The parsed JSON is
// stored in the value pointed to by v
func GetJSON(url string, v interface{}) error {
	body, err := Get(url)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
