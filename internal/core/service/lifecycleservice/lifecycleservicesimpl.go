package lifecycleservice

import (
	"github.com/rameshpolishetti/mlca/internal/core/common/config"
	"github.com/rameshpolishetti/mlca/internal/core/component"
	"github.com/rameshpolishetti/mlca/internal/core/service"
	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("lifecycle-service")

// LifeCycleServices
type LifeCycleServices interface {
	CheckState() bool
}

// LifeCycleServicesImpl LifeCycleServiceImpl
type LifeCycleServicesImpl struct {
	containerDaemon config.ContainerDaemon
	managedServices map[string]LifeCycleService
	regService      *service.RegistryProxy
}

// NewLifeCycleServices creates new LifeCycleServiceImpl
func NewLifeCycleServices(cDaemon config.ContainerDaemon) LifeCycleServices {

	// load managed services map
	mServices := make(map[string]LifeCycleService)
	rService := service.NewRegistryProxyService(cDaemon)

	for _, c := range cDaemon.Components {
		var mc component.Component
		if c.Type == "microgateway" {
			mc = component.NewMicrogatewayComponent(c.Name)
		} else {
			log.Panicf("managed component of type %s not found", c.Type)
		}
		mServices[c.Name] = NewLifeCycleService(mc, rService)
	}

	lcServicesImpl := &LifeCycleServicesImpl{
		containerDaemon: cDaemon,
		managedServices: mServices,
		regService:      rService,
	}

	return lcServicesImpl
}

// CheckState check managed component state
func (lcServicesImpl *LifeCycleServicesImpl) CheckState() bool {
	result := false
	if !lcServicesImpl.regService.IsReady() {
		log.Info("Registry is not ready")
		return result
	}
	for _, mService := range lcServicesImpl.managedServices {
		if mService.CheckState() {
			result = true
		}
	}
	return result
}
