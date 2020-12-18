package restutil

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Get performs a GET request and return the body content
func Get(url string, headers http.Header) ([]byte, error) {
	c := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("unable to create request, " + err.Error())
	}

	req.Header = headers

	response, err := c.Do(req)
	if err != nil {
		return nil, errors.New("unable to perform request, " + err.Error())
	}

	reader := response.Body
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.New("unable to read the response body, " + err.Error())
	}

	return body, nil
}
