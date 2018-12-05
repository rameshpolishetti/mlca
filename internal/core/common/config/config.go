package config

// ContainerDaemon container configuration
type ContainerDaemon struct {
	Name              string             `json:"name"`
	Domain            string             `json:"domain"`
	Cluster           string             `json:"cluster"`
	Zone              string             `json:"zone"`
	Node              string             `json:"node"`
	Type              string             `json:"type"`
	Qualifier         string             `json:"qualifier"`
	Inboxes           map[string]string  `json:"inboxes"`
	TransportSettings TransportSettings  `json:"transportSettings"`
	Components        []ManagedComponent `json:"components"`

	IP string
}

// ManagedComponent component configuration
type ManagedComponent struct {
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	Qualifier         string            `json:"qualifier"`
	Script            string            `json:"script"`
	Service           string            `json:"service"`
	Factory           string            `json:"factory"`
	ContainerInstance ContainerInstance `json:"container"`
}

// TransportSettings transport configuration
type TransportSettings struct {
	Scheme string `json:"scheme"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
}

// ContainerInstance container instance configuration
type ContainerInstance struct {
	Domain  string `json:"domain"`
	Cluster string `json:"cluster"`
	Zone    string `json:"zone"`
	Node    string `json:"node"`

	DomainID  string `json:"domainId"`
	ClusterID string `json:"clusterId"`
	ZoneID    string `json:"zoneId"`
	NodeID    string `json:"nodeId"`

	IP              string `json:"ip"`
	RegistryContext string `json:"registryContext"`
}

func (mc *ManagedComponent) Clone(copyFrom ManagedComponent) {
	mc.Name = copyFrom.Name
	mc.Type = copyFrom.Type
	mc.Qualifier = copyFrom.Qualifier
	mc.Script = copyFrom.Script
	mc.Service = copyFrom.Service
	mc.Factory = copyFrom.Factory
	// TODO
	// mc.ContainerInstance = copyFrom.ContainerInstance
}
