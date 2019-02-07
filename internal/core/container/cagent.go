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
	"github.com/rameshpolishetti/mlca/internal/core/service"
	"github.com/rameshpolishetti/mlca/internal/core/service/lifecycleservice"
)

var log = logger.GetLogger("cagent")

const (
	HeartBeatInterval = 2000 * time.Millisecond
)

// ContainerAgent container agent
type ContainerAgent struct {
	containerDaemon   config.ContainerDaemon
	RegService        *service.RegistryProxy
	LifecycleServices lifecycleservice.LifeCycleServices
}

// NewContainerAgent creates new container agent
func NewContainerAgent(cDaemon config.ContainerDaemon) *ContainerAgent {

	a := &ContainerAgent{
		containerDaemon: cDaemon,
	}

	// Init registry proxy service
	a.RegService = service.NewRegistryProxyService(cDaemon)

	// load managed components

	// init lifecycle services
	a.LifecycleServices = lifecycleservice.NewLifeCycleServices(cDaemon)

	return a
}

// Initialize initializes container agent
func (ca *ContainerAgent) Initialize() {
	// it it lifecycle service
}

// Start starts container agent
func (ca *ContainerAgent) Start() {

	// start http server
	router := mux.NewRouter()
	pathStatus := fmt.Sprintf("/%s/status", ca.containerDaemon.Name)
	router.HandleFunc(pathStatus, ca.getStatus).Methods("GET")
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", ca.containerDaemon.TransportSettings.Port),
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
				// if ca.switchState() {
				// 	// update registry with new state
				// 	ca.RegService.UpdateStatus(ca.FSM.Current())
				// } else {
				// 	log.Debugln("no state transition")
				// }
				if !ca.LifecycleServices.CheckState() {
					log.Error("CheckState FAIL")
				}

			case <-signalChan:
				log.Infoln("Received os interrupt")

				// stop container agent
				log.Infoln("Stop container agent")
				hearBeatTimer.Stop()

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

// ModelCA model container agent
type ModelCA struct {
	Name   string `json:"name"`
	Status string `json:status`
}

// REST API
func (ca *ContainerAgent) getStatus(w http.ResponseWriter, r *http.Request) {
	mca := &ModelCA{
		Name:   ca.containerDaemon.Name,
		Status: "UNKNOWN",
	}

	json.NewEncoder(w).Encode(mca)
}
