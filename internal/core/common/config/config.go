package config

// ContainerConfig container configuration
type ContainerConfig struct {
	Name              string            `json:"name"`
	Domain            string            `json:"domain"`
	Cluster           string            `json:"cluster"`
	Zone              string            `json:"zone"`
	Node              string            `json:"node"`
	Type              string            `json:"type"`
	Qualifier         string            `json:"qualifier"`
	Inboxes           map[string]string `json:"inboxes"`
	TransportSettings TransportSettings `json:"transportSettings"`
	Components        []Component       `json:"components"`

	IP string
}

// Component component configuration
type Component struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Qualifier string `json:"qualifier"`
	Script    string `json:"script"`
	Service   string `json:"service"`
	Factory   string `json:"factory"`
}

// TransportSettings transport configuration
type TransportSettings struct {
	Scheme string `json:"scheme"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
}
