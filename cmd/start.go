package cmd

import (
	"github.com/rameshpolishetti/mlca/logger"

	"github.com/rameshpolishetti/mlca/internal/cagent"
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
	containerName := viper.GetString("name")
	registry := viper.GetString("inboxes.registry")
	caPort := viper.GetInt("transportSettings.port")

	log.Infof("Start the container agent with the container - %s", containerName)

	ca := cagent.NewContainerAgent(containerName, registry, caPort)
	ca.Start()
}
