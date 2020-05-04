package restutil

import (
	"encoding/csv"
	"errors"
	"net/http"
)

func GetCSV(url string) ([][]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.New("unable to perform request, " + err.Error())
	}

	reader := response.Body
	defer reader.Close()

	return csv.NewReader(reader).ReadAll()
}
