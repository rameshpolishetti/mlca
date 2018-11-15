package container

// Config container configuration
type Config struct {
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
}

// Component component configuration
type Component map[string]string

// TransportSettings transport configuration
type TransportSettings struct {
	Scheme string `json:"scheme"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
}
