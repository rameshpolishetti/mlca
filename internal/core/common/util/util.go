package util

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("util")

// LookupHostIP looks up host IP
func LookupHostIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

// RunScript runs script
func RunScript(script string) {
	log.Infof("Running the script [%s]", script)
	scriptTokens := strings.Split(script, " ")
	cmd := exec.Command(scriptTokens[0], scriptTokens[1:]...)
	err := cmd.Start()
	if err != nil {
		log.Fatal("Not able to run the script with error - ", err)
	}
}
