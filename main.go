package main

import (
	"github.com/rameshpolishetti/mlca/cmd"
	"github.com/rameshpolishetti/mlca/logger"
)

var log = logger.GetLogger("main")

func main() {
	cmd.Execute()
}
