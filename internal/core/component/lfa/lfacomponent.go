package lfa

import (
	"github.com/rameshpolishetti/mlca/internal/core/common/config"
	"github.com/rameshpolishetti/mlca/internal/core/common/util"
	"github.com/rameshpolishetti/mlca/internal/core/component"
	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("lfa")

// LFAComponent holds log forwarding agent component
type LFAComponent struct {
	// Name string
	config.ManagedComponent
}

// NewLFAComponent creates new LFAComponent
func NewLFAComponent(mc config.ManagedComponent) component.Component {
	log.Infoln("init")
	lfaComponent := &LFAComponent{}
	lfaComponent.Clone(mc)

	return lfaComponent
}

func (lfac *LFAComponent) Bootup() bool {
	log.Infoln("Bootup")
	return true
}

func (lfac *LFAComponent) BuildConfiguration() bool {
	log.Infoln("BuildConfiguration")
	return true
}

func (lfac *LFAComponent) LaunchComponent() bool {
	log.Infoln("LaunchComponent")
	return true
}

func (lfac *LFAComponent) PrepareForActive() bool {
	log.Infoln("PrepareForActive")
	// run script
	util.RunScript(lfac.Script)
	return true
}

func (lfac *LFAComponent) WatchComponent() bool {
	log.Infoln("WatchComponent")
	return true
}
