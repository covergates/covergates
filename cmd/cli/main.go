package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	log.Info("hi")
}
