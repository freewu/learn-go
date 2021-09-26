package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears") // time="2021-09-26T15:40:15+08:00" level=info msg="A walrus appears" animal=walrus
}