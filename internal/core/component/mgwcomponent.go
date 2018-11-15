package component

import (
	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("tm")

// MicrogatewayComponent holds Microgateway component
type MicrogatewayComponent struct {
	Name string
}

// NewMicrogatewayComponent creates new MicrogatewayComponent component
func NewMicrogatewayComponent(name string) Component {
	log.Infoln("init")
	mgwComponent := &MicrogatewayComponent{
		Name: name,
	}

	return mgwComponent
}

func (mgwc *MicrogatewayComponent) Bootup() bool {
	log.Infoln("Bootup")
	return true
}

func (mgwc *MicrogatewayComponent) BuildConfiguration() bool {
	log.Infoln("BuildConfiguration")
	return true
}

func (mgwc *MicrogatewayComponent) LaunchComponent() bool {
	log.Infoln("LaunchComponent")
	return true
}

func (mgwc *MicrogatewayComponent) PrepareForActive() bool {
	log.Infoln("PrepareForActive")
	return true
}

func (mgwc *MicrogatewayComponent) WatchComponent() bool {
	log.Infoln("WatchComponent")
	return true
}
