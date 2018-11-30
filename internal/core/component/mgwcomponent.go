package component

import (
	"github.com/project-flogo/contrib/activity/rest"
	trigger "github.com/project-flogo/contrib/trigger/rest"
	"github.com/project-flogo/core/api"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/microgateway"
	microapi "github.com/project-flogo/microgateway/api"
	"github.com/rameshpolishetti/mlca/internal/core/common/config"
	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("tm")

// MicrogatewayComponent holds Microgateway component
type MicrogatewayComponent struct {
	// Name string
	config.ManagedComponent
}

// NewMicrogatewayComponent creates new MicrogatewayComponent component
func NewMicrogatewayComponent(name string) Component {
	log.Infoln("init")
	mgwComponent := &MicrogatewayComponent{
	// Name: name,
	}
	mgwComponent.Name = name

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
	// startMicrogateway()
	return true
}

func (mgwc *MicrogatewayComponent) WatchComponent() bool {
	log.Infoln("WatchComponent")
	return true
}

// microgateway
func startMicrogateway() {
	app := api.NewApp()

	gateway := microapi.New("Pets")
	service := gateway.NewService("PetStorePets", &rest.Activity{})
	service.SetDescription("Get pets by ID from the petstore")
	service.AddSetting("uri", "http://petstore.swagger.io/v2/pet/:petId")
	service.AddSetting("method", "GET")
	step := gateway.NewStep(service)
	step.AddInput("pathParams", "=$.payload.pathParams")
	response := gateway.NewResponse(false)
	response.SetCode(200)
	response.SetData("=$.PetStorePets.outputs.data")
	settings, err := gateway.AddResource(app)
	if err != nil {
		panic(err)
	}

	trg := app.NewTrigger(&trigger.Trigger{}, &trigger.Settings{Port: 9096})
	handler, err := trg.NewHandler(&trigger.HandlerSettings{
		Method: "GET",
		Path:   "/pets/:petId",
	})
	if err != nil {
		panic(err)
	}

	_, err = handler.NewAction(&microgateway.Action{}, settings)
	if err != nil {
		panic(err)
	}

	e, err := api.NewEngine(app)
	if err != nil {
		panic(err)
	}
	go engine.RunEngine(e)
}
