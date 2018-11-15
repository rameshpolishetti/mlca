package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("registry-service")

type RegistryProxy struct {
	name       string
	baseURL    *url.URL
	httpClient *http.Client

	// registry status
	isReady bool

	// cluster info
	tmgcId    string
	zoneId    string
	clusterId string
}

func NewRegistryProxyService(name string, registry string) *RegistryProxy {
	rp := &RegistryProxy{
		name: name,
		baseURL: &url.URL{
			Scheme: "http",
			Host:   registry,
			Path:   "registry/rest/v1",
		},
		httpClient: &http.Client{},
	}
	return rp
}

func (rp *RegistryProxy) Register() bool {
	// check whether the registry is ready
	if !rp.IsReady() {
		log.Infoln("Registry is not ready")
		return false
	}

	registryPath := rp.baseURL.String() + "/clusters/Mashery/zones/Local/" + rp.name
	log.Infoln("POST request to: ", registryPath)
	/**
	 * create registry payload
	 * {
	 * "name": "proxy-node1",
	 * "host": "10.1.2.3",
	 * "agentPort": 1234,
	 * "status": "registering",
	 * ... plus container specific arguments
	 * }
	 */
	payloadMap := map[string]interface{}{
		"name":      "microgateway",
		"host":      "10.97.90.65",
		"agentPort": 21780,
		"status":    "registering",
	}
	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		log.Errorln(err)
	}
	log.Debugf("request payload: %s", payloadBytes)

	res, err := rp.httpClient.Post(registryPath, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Errorln(err)
		return false
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorln(err)
		return false
	}
	log.Infof("Response form registry: %s \n", body)
	/* sample response
	{
		"registrationTime" : "11-317-18 15:46:53.914+0530",
		"tmgcId" : "9e86528a-f7b1-415a-bbdc-048185395c64",
		"host" : "10.97.90.65",
		"name" : "microgateway",
		"zoneId" : "9b8e3d94-d32a-4981-874b-c8e0697afe13",
		"clusterId" : "5057c530-bbae-4210-8c46-62fc02581618",
		"agentPort" : 21780,
		"status" : "registered"
	}
	*/
	type RegistryResp struct {
		TmgcId    string `json:tmgcId`
		ZoneId    string `json:zoneId`
		ClusterId string `json:clusterId`
		Status    string `json:status`
	}
	respObj := &RegistryResp{}
	err = json.Unmarshal([]byte(body), respObj)
	if err != nil {
		log.Errorln(err)
		return false
	}

	if respObj.Status == "registered" {
		rp.tmgcId = respObj.TmgcId
		rp.zoneId = respObj.ZoneId
		rp.clusterId = respObj.ClusterId
		return true
	}
	return false
}

func (rp *RegistryProxy) IsReady() bool {
	if rp.isReady {
		return true
	}

	statusPath := rp.baseURL.String() + "/status"

	res, err := rp.httpClient.Get(statusPath)
	if err != nil {
		log.Errorln(err)
		return false
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorln(err)
		return false
	}

	type RegistryResp struct {
		Status string `json:"status"`
	}
	respObj := &RegistryResp{}

	err = json.Unmarshal([]byte(body), respObj)
	if err != nil {
		log.Errorln(err)
		return false
	}

	if respObj.Status == "REGISTRY_READY" {
		rp.isReady = true
	}
	return rp.isReady
}

func (rp *RegistryProxy) UpdateStatus(status string) bool {
	// check whether the registry is ready
	if !rp.IsReady() {
		log.Infoln("Registry is not ready")
		return false
	}
	/*
		payload: {"status":"UNSATISFIED"}
		path: /clusters/<>/zones/<>/containerName/<>/status
		sample response:
			{
				"updatedTime" : "11-318-18 19:16:46.998+0530",
				"tmgcId" : "d530176c-d85c-4160-b18f-f46377f104bf",
				"tmgcType" : "microgateway",
				"zoneId" : "ba7eb03c-8616-4b0f-a379-ec6bc2e5ab33",
				"clusterId" : "bf9cd2a6-32ec-4d34-acbb-429449e6af88",
				"message" : "Updated status of the tmgc",
				"status" : "UNSATISFIED"
			}
	*/
	statusPath := rp.baseURL.String() +
		"/clusters/" + rp.clusterId +
		"/zones/" + rp.zoneId +
		"/" + rp.name + "/" + rp.tmgcId +
		"/status"
	log.Infoln("PUT request to: ", statusPath)

	payloadMap := map[string]interface{}{
		"status": status,
	}
	payloadBytes, err := json.Marshal(payloadMap)
	if err != nil {
		log.Errorln(err)
	}
	log.Debugf("request payload: %s", payloadBytes)

	req, err := http.NewRequest(http.MethodPut, statusPath, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Errorln(err)
		return false
	}
	// res, err := rp.httpClient.Post(statusPath, "application/json", bytes.NewBuffer(payloadBytes))
	res, err := rp.httpClient.Do(req)
	if err != nil {
		log.Errorln(err)
		return false
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorln(err)
		return false
	}
	log.Infof("Updated status in registry to %s - Response from registry: %s", status, body)

	return true
}
