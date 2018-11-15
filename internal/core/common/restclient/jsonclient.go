package restclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("jsonclient")

// Get performs http GET
func Get(path string) ([]byte, error) {
	httpClient := getHTTPClient()
	log.Debugf("GET request to %s", path)
	res, err := httpClient.Get(path)
	if err != nil {
		log.Errorf("GET request to %s failed. Reason: %s", path, err)
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("GET request to %s failed. Reason: %s", path, err)
		return nil, err
	}

	return resBody, nil
}

// Post performs http POST
func Post(path string, payloadMap map[string]interface{}) ([]byte, error) {
	log.Debugf("POST request to %s", path)
	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		log.Errorf("POST request to %s failed. Reason: %s", path, err)
		return nil, err
	}
	log.Debugf("payload: %s", payloadBytes)

	httpClient := getHTTPClient()
	res, err := httpClient.Post(path, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Errorf("POST request to %s failed. Reason: %s", path, err)
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("POST request to %s failed. Reason: %s", path, err)
		return nil, err
	}

	return resBody, nil
}

// Put performs http PUT
func Put(path string, payloadMap map[string]interface{}) ([]byte, error) {
	log.Debugf("PUT request to %s", path)
	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		log.Errorf("PUT request to %s failed. Reason: %s", path, err)
	}
	log.Debugf("payload: %s", payloadBytes)

	req, err := http.NewRequest(http.MethodPut, path, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Errorf("PUT request to %s failed. Reason: %s", path, err)
		return nil, err
	}

	httpClient := getHTTPClient()
	res, err := httpClient.Do(req)
	if err != nil {
		log.Errorf("PUT request to %s failed. Reason: %s", path, err)
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("PUT request to %s failed. Reason: %s", path, err)
		return nil, err
	}

	return resBody, nil
}

func getHTTPClient() *http.Client {
	httpClient := &http.Client{}
	return httpClient
}
