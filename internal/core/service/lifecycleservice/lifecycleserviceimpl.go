package lifecycleservice

import (
	"github.com/looplab/fsm"
	"github.com/rameshpolishetti/mlca/internal/core/component"
	"github.com/rameshpolishetti/mlca/internal/core/service"
)

// LifeCycleService LifeCycleService
type LifeCycleService interface {
	CheckState() bool
}

// LifeCycleServiceImpl LifeCycleServiceImpl
type LifeCycleServiceImpl struct {
	FSM        *fsm.FSM
	mComponent component.Component
	regService *service.RegistryProxy
}

// NewLifeCycleService New
func NewLifeCycleService(mc component.Component, rService *service.RegistryProxy) LifeCycleService {
	lcServiceImpl := &LifeCycleServiceImpl{
		mComponent: mc,
		regService: rService,
	}

	/*
	* UNKNOWN	initialize()	bootup()
	* UNSATISFIED	resolveDependencies()	buildConfiguration()
	* RESOLVED	activate()	launchComponent()
	* STANDBY	standby()	prepareForActive()
	* ACTIVE	monitor()	watchComponent()
	* RELOAD	reload()
	* RECYCLE	waitingForDependencies()
	* DISABLED	deavtivate()
	 */

	lcServiceImpl.FSM = fsm.NewFSM(
		"UNKNOWN",
		fsm.Events{
			{Name: "initialize", Src: []string{"UNKNOWN"}, Dst: "UNSATISFIED"},
			{Name: "resolveDependencies", Src: []string{"UNSATISFIED"}, Dst: "RESOLVED"},
			{Name: "activate", Src: []string{"RESOLVED"}, Dst: "STANDBY"},
			{Name: "standby", Src: []string{"STANDBY"}, Dst: "ACTIVE"},
			{Name: "monitor", Src: []string{"ACTIVE"}, Dst: "ACTIVE"},
			{Name: "deavtivate", Src: []string{"ACTIVE"}, Dst: "UNKNOWN"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { lcServiceImpl.enterState(e) },
		},
	)
	return lcServiceImpl
}

// CheckState CheckState
func (lcServiceImpl *LifeCycleServiceImpl) enterState(e *fsm.Event) {
	log.Debugf("%s -> %s", e.Src, e.Dst)
}

// CheckState CheckState
func (lcServiceImpl *LifeCycleServiceImpl) CheckState() bool {
	if lcServiceImpl.switchState() {
		// update registry status
		lcServiceImpl.regService.UpdateStatus(lcServiceImpl.FSM.Current())
		return true
	}
	return false
}

func (lcServiceImpl *LifeCycleServiceImpl) switchState() bool {
	if lcServiceImpl.FSM.Is("UNKNOWN") && lcServiceImpl.initialize() {
		return true
	}

	if lcServiceImpl.FSM.Is("UNSATISFIED") && lcServiceImpl.resolveDependencies() {
		return true
	}

	if lcServiceImpl.FSM.Is("RESOLVED") && lcServiceImpl.activate() {
		return true
	}

	if lcServiceImpl.FSM.Is("STANDBY") && lcServiceImpl.standby() {
		return true
	}

	if lcServiceImpl.FSM.Is("ACTIVE") && lcServiceImpl.monitor() {
		return true
	}

	return false
}

func (lcServiceImpl *LifeCycleServiceImpl) initialize() bool {
	// init
	// bootup() -> register -> ConfigurationRegistryService -> register()

	// bootup component
	if !lcServiceImpl.mComponent.Bootup() {
		return false
	}

	// register
	if lcServiceImpl.regService.Register() {
		log.Infoln("Registration SUCCESS")

		// update state
		err := lcServiceImpl.FSM.Event("initialize")
		if err != nil {
			log.Errorln(err)
			return false
		}
		return true
	}
	log.Infoln("Registration FAIL")

	return false
}

func (lcServiceImpl *LifeCycleServiceImpl) resolveDependencies() bool {
	// resolve
	if !lcServiceImpl.mComponent.BuildConfiguration() {
		return false
	}

	// update state
	err := lcServiceImpl.FSM.Event("resolveDependencies")
	if err != nil {
		log.Errorln(err)
		return false
	}
	return true
}

func (lcServiceImpl *LifeCycleServiceImpl) activate() bool {
	// activate
	if !lcServiceImpl.mComponent.LaunchComponent() {
		return false
	}
	// update state
	err := lcServiceImpl.FSM.Event("activate")
	if err != nil {
		log.Errorln(err)
		return false
	}

	return true
}

func (lcServiceImpl *LifeCycleServiceImpl) standby() bool {
	// standby
	if !lcServiceImpl.mComponent.PrepareForActive() {
		return false
	}
	// update state
	err := lcServiceImpl.FSM.Event("standby")
	if err != nil {
		log.Errorln(err)
		return false
	}
	return true
}

func (lcServiceImpl *LifeCycleServiceImpl) monitor() bool {
	// monitor
	if !lcServiceImpl.mComponent.WatchComponent() {
		return false
	}
	// update state
	err := lcServiceImpl.FSM.Event("monitor")
	if err != nil && err.Error() != "no transition" {
		log.Errorln(err)
		lcServiceImpl.FSM.SetState("RESOLVED")
		return false
	}
	log.Infof("[monitor] Current state: %s", lcServiceImpl.FSM.Current())
	return true
}

func (lcServiceImpl *LifeCycleServiceImpl) deavtivate() bool {
	// deavtivate
	err := lcServiceImpl.FSM.Event("deavtivate")
	if err != nil {
		log.Errorln(err)
		return false
	}
	return true
}
