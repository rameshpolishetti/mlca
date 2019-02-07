package cmd

import (
	"encoding/json"
	"os"

	"github.com/rameshpolishetti/mlca/internal/core/common/config"
	"github.com/rameshpolishetti/mlca/internal/core/common/util"
	"github.com/rameshpolishetti/mlca/internal/core/container"
	"github.com/rameshpolishetti/mlca/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logger.GetLogger("cmd")
var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	startCmd.Flags().StringVarP(&cfgFile, "config", "c", "config.json", "configuration file")
	// startCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(startCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		log.Panicln("config file is required")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Panicln("Can't read config:", err)
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run containerized application",
	Long:  `Start runs containerized application`,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {

	// load container configuration
	var cConfig config.ContainerDaemon
	err := viper.Unmarshal(&cConfig)
	if err != nil {
		log.Errorf("unable to load container configuration")
	}

	// load host details
	cConfig.IP = util.LookupHostIP()
	hName, err := os.Hostname()
	if err != nil {
		log.Errorf("unable to load host details")
	}
	cConfig.Node = hName

	cfgString, _ := json.MarshalIndent(cConfig, "", " ")
	log.Infof("Start the container [%s] with configuration: %s", cConfig.Name, cfgString)

	ca := container.NewContainerAgent(cConfig)
	ca.Initialize()
	ca.Start()
}
