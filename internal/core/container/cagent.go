package container

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rameshpolishetti/mlca/logger"

	"github.com/gorilla/mux"
	"github.com/rameshpolishetti/mlca/internal/core/common/config"
	"github.com/rameshpolishetti/mlca/internal/core/component"
	"github.com/rameshpolishetti/mlca/internal/core/service"

	"github.com/looplab/fsm"
)

var log = logger.GetLogger("cagent")

const (
	HeartBeatInterval = 2000 * time.Millisecond
)

// ContainerAgent container agent
type ContainerAgent struct {
	Config     config.ContainerConfig
	FSM        *fsm.FSM
	Component  component.Component
	RegService *service.RegistryProxy
}

// NewContainerAgent creates new container agent
func NewContainerAgent(cConfig config.ContainerConfig) *ContainerAgent {

	a := &ContainerAgent{
		Config: cConfig,
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

	a.FSM = fsm.NewFSM(
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
			"enter_state": func(e *fsm.Event) { a.enterState(e) },
		},
	)

	// fmt.Println(fsm.Visualize(a.FSM))

	// Init component
	a.Component = component.NewMicrogatewayComponent(cConfig.Name)

	// Init registry proxy service
	a.RegService = service.NewRegistryProxyService(cConfig)

	return a
}

func (ca *ContainerAgent) enterState(e *fsm.Event) {
	log.Infof("%s -> %s", e.Src, e.Dst)
}

// Start starts container agent
func (ca *ContainerAgent) Start() {

	// start http server
	router := mux.NewRouter()
	pathStatus := fmt.Sprintf("/%s/status", ca.Config.Name)
	router.HandleFunc(pathStatus, ca.getStatus).Methods("GET")
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", ca.Config.TransportSettings.Port),
		Handler: router,
	}
	go func() {
		log.Infoln("Start http server")
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// os signal channel (ctrl+c)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	// heart beat timer
	hearBeatTimer := time.NewTicker(HeartBeatInterval)

	// exit channel
	exitChan := make(chan int)

	// start lifecycle
	go func() {
		for {
			select {
			case <-hearBeatTimer.C:
				// switch state
				if ca.switchState() {
					// update registry with new state
					ca.RegService.UpdateStatus(ca.FSM.Current())
				} else {
					log.Debugln("no state transition")
				}

			case <-signalChan:
				log.Infoln("Received os interrupt")

				// stop container agent
				log.Infoln("Stop container agent")
				hearBeatTimer.Stop()
				ca.deavtivate()

				// shutdown http server
				log.Infoln("Shutting down http server")
				httpServer.Shutdown(context.Background())

				// exit
				exitChan <- 0
			}
		}
	}()

	code := <-exitChan
	os.Exit(code)
}

func (ca *ContainerAgent) switchState() bool {
	if ca.FSM.Is("UNKNOWN") && ca.initialize() {
		return true
	}

	if ca.FSM.Is("UNSATISFIED") && ca.resolveDependencies() {
		return true
	}

	if ca.FSM.Is("RESOLVED") && ca.activate() {
		return true
	}

	if ca.FSM.Is("STANDBY") && ca.standby() {
		return true
	}

	if ca.FSM.Is("ACTIVE") && ca.monitor() {
		return true
	}

	return false
}

func (ca *ContainerAgent) initialize() bool {
	// init
	// bootup() -> register -> ConfigurationRegistryService -> register()

	// bootup component
	if !ca.Component.Bootup() {
		return false
	}

	// register
	if ca.RegService.Register() {
		log.Infoln("Registration SUCCESS")

		// update state
		err := ca.FSM.Event("initialize")
		if err != nil {
			log.Errorln(err)
			return false
		}
		return true
	}
	log.Infoln("Registration FAIL")

	return false
}

func (ca *ContainerAgent) resolveDependencies() bool {
	// resolve
	if !ca.Component.BuildConfiguration() {
		return false
	}

	// update state
	err := ca.FSM.Event("resolveDependencies")
	if err != nil {
		log.Errorln(err)
		return false
	}
	return true
}

func (ca *ContainerAgent) activate() bool {
	// activate
	if !ca.Component.LaunchComponent() {
		return false
	}
	// update state
	err := ca.FSM.Event("activate")
	if err != nil {
		log.Errorln(err)
		return false
	}

	return true
}

func (ca *ContainerAgent) standby() bool {
	// standby
	if !ca.Component.PrepareForActive() {
		return false
	}
	// update state
	err := ca.FSM.Event("standby")
	if err != nil {
		log.Errorln(err)
		return false
	}
	return true
}

func (ca *ContainerAgent) monitor() bool {
	// monitor
	if !ca.Component.WatchComponent() {
		return false
	}
	// update state
	err := ca.FSM.Event("monitor")
	if err != nil && err.Error() != "no transition" {
		log.Errorln(err)
		ca.FSM.SetState("RESOLVED")
		return false
	}
	log.Infof("[monitor] Current state: %s", ca.FSM.Current())
	return true
}

func (ca *ContainerAgent) deavtivate() bool {
	// deavtivate
	err := ca.FSM.Event("deavtivate")
	if err != nil {
		log.Errorln(err)
		return false
	}
	return true
}

// ModelCA model container agent
type ModelCA struct {
	Name   string `json:"name"`
	Status string `json:status`
}

// REST API
func (ca *ContainerAgent) getStatus(w http.ResponseWriter, r *http.Request) {
	mca := &ModelCA{
		Name:   ca.Config.Name,
		Status: ca.FSM.Current(),
	}

	json.NewEncoder(w).Encode(mca)
}
