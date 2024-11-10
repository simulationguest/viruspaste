package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Response struct {
	Data struct {
		Attributes struct {
			LastAnalysisStats struct {
				Malicious int `json:"malicious"`
			} `json:"last_analysis_stats"`
		} `json:"attributes"`
	} `json:"data"`
}

var apiKey = os.Getenv("API_KEY")

func checkFile(hash string) (bool, error) {
	url := "https://www.virustotal.com/api/v3/files/" + hash

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-apikey", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	// if res.StatusCode != http.StatusOK {
	// 	return false, errors.New("bad request")
	// }

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	response := &Response{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return false, err
	}

	return response.Data.Attributes.LastAnalysisStats.Malicious == 0, nil
}
