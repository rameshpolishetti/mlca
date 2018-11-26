package restclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("jsonclient")

// JSONClient json client utility for http client operations like GET, POST, PUT, etc.
type JSONClient struct {
	inbox   string
	context string
}

// New creates new JSONClient
func New(i, c string) *JSONClient {
	jsonClient := &JSONClient{
		inbox:   i,
		context: c,
	}
	return jsonClient
}

func (jsonClient *JSONClient) getRequestURL(path string) (string, error) {
	u, err := url.Parse(jsonClient.inbox + jsonClient.context + path)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// Get performs http GET
func (jsonClient *JSONClient) Get(path string) ([]byte, error) {
	requestURL, err := jsonClient.getRequestURL(path)
	if err != nil {
		return nil, err
	}

	httpClient := getHTTPClient()
	log.Debugf("GET request to %s", requestURL)
	res, err := httpClient.Get(requestURL)
	if err != nil {
		log.Errorf("GET request to %s failed. Reason: %s", requestURL, err)
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("GET request to %s failed. Reason: %s", requestURL, err)
		return nil, err
	}

	return resBody, nil
}

// Post performs http POST
func (jsonClient *JSONClient) Post(path string, payloadMap map[string]interface{}) ([]byte, error) {
	requestURL, err := jsonClient.getRequestURL(path)
	if err != nil {
		return nil, err
	}

	log.Debugf("POST request to %s", requestURL)
	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		log.Errorf("POST request to %s failed. Reason: %s", requestURL, err)
		return nil, err
	}
	log.Debugf("payload: %s", payloadBytes)

	httpClient := getHTTPClient()
	res, err := httpClient.Post(requestURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Errorf("POST request to %s failed. Reason: %s", requestURL, err)
		return nil, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("POST request to %s failed. Reason: %s", requestURL, err)
		return nil, err
	}

	return resBody, nil
}

// Put performs http PUT
func (jsonClient *JSONClient) Put(path string, payloadMap map[string]interface{}) ([]byte, error) {
	requestURL, err := jsonClient.getRequestURL(path)
	if err != nil {
		return nil, err
	}

	log.Debugf("PUT request to %s", requestURL)
	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		log.Errorf("PUT request to %s failed. Reason: %s", requestURL, err)
	}
	log.Debugf("payload: %s", payloadBytes)

	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Errorf("PUT request to %s failed. Reason: %s", requestURL, err)
		return nil, err
	}

	httpClient := getHTTPClient()
	res, err := httpClient.Do(req)
	if err != nil {
		log.Errorf("PUT request to %s failed. Reason: %s", requestURL, err)
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("PUT request to %s failed. Reason: %s", requestURL, err)
		return nil, err
	}

	return resBody, nil
}

func getHTTPClient() *http.Client {
	httpClient := &http.Client{}
	return httpClient
}
