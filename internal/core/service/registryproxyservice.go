package service

import (
	"encoding/json"

	"github.com/rameshpolishetti/mlca/internal/core/common/config"
	jsonclient "github.com/rameshpolishetti/mlca/internal/core/common/restclient"
	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("registry-service")

// RegistryProxy rigistry
type RegistryProxy struct {
	cConfig    config.ContainerDaemon
	jsonClient *jsonclient.JSONClient

	// registry status
	isReady bool

	// cluster info
	tmgcId    string
	zoneId    string
	clusterId string
}

// NewRegistryProxyService creates new registry proxy
func NewRegistryProxyService(cCfg config.ContainerDaemon) *RegistryProxy {
	registry := cCfg.Inboxes["registry"]
	registryContext := "/registry/rest/v1"
	rp := &RegistryProxy{
		cConfig:    cCfg,
		jsonClient: jsonclient.New(registry, registryContext),
	}
	return rp
}

// Register rigister with rigistry
func (rp *RegistryProxy) Register() bool {
	// check whether the registry is ready
	if !rp.IsReady() {
		log.Infoln("Registry is not ready")
		return false
	}

	registerPath := "/clusters/" + rp.cConfig.Cluster + "/zones/" + rp.cConfig.Zone + "/" + rp.cConfig.ComponentType
	log.Infoln("Registering")
	/**
	 * create registry payload
	 * {
	 * "name": "proxy-node1",
	 * "host": "10.1.2.3",
	 * "port": 9096,
	 * "agentPort": 1234,
	 * "status": "registering",
	 * ... plus container specific arguments
	 * }
	 */
	payloadMap := map[string]interface{}{
		"name":      rp.cConfig.Name,
		"host":      rp.cConfig.IP,
		"port":      rp.cConfig.Port,
		"agentPort": rp.cConfig.TransportSettings.Port,
		"status":    "registering",
	}

	res, err := rp.jsonClient.Post(registerPath, payloadMap)
	if err != nil {
		return false
	}
	log.Infof("Response from registry: %s \n", res)
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
	err = json.Unmarshal(res, respObj)
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

// IsReady return whether registry is ready
func (rp *RegistryProxy) IsReady() bool {
	if rp.isReady {
		return true
	}

	statusPath := "/status"

	body, err := rp.jsonClient.Get(statusPath)
	if err != nil {
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

// UpdateStatus updates status with registry
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
	updateStatusPath := "/clusters/" + rp.clusterId +
		"/zones/" + rp.zoneId +
		"/" + rp.cConfig.ComponentType + "/" + rp.tmgcId +
		"/status"
	log.Infoln("PUT request to: ", updateStatusPath)

	payloadMap := map[string]interface{}{
		"status": status,
	}

	res, err := rp.jsonClient.Put(updateStatusPath, payloadMap)
	if err != nil {
		return false
	}
	log.Infof("Updated status in registry to %s - Response from registry: %s", status, res)

	return true
}
