package mgw

import (
	"github.com/rameshpolishetti/mlca/internal/core/common/config"
	"github.com/rameshpolishetti/mlca/internal/core/common/util"
	"github.com/rameshpolishetti/mlca/internal/core/component"
	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("tm")

// MicrogatewayComponent holds Microgateway component
type MicrogatewayComponent struct {
	// Name string
	config.ManagedComponent
}

// NewMicrogatewayComponent creates new MicrogatewayComponent component
func NewMicrogatewayComponent(mc config.ManagedComponent) component.Component {
	log.Infoln("init")
	mgwComponent := &MicrogatewayComponent{
	// Name: name,
	}
	mgwComponent.Clone(mc)

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
	// run script
	util.RunScript(mgwc.Script)
	return true
}

func (mgwc *MicrogatewayComponent) WatchComponent() bool {
	log.Infoln("WatchComponent")
	return true
}
