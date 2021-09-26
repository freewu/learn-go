package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	//log.SetFormatter(&log.TextFormatter{}) // default

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// log.SetOutput(os.Stderr) // default

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel) // trace / debug / info / warn / error / fatal / panic 设置 warn 级别以上输出
}

func main() {
	log.Trace("Trace message") //
	log.Debug("Debug message") //
	log.Info("Info message") //
	log.Warn("Warn message") // {"level":"warning","msg":"Warn message","time":"2021-09-26T16:08:05+08:00"}
	log.Error("Error message") // {"level":"error","msg":"Error message","time":"2021-09-26T16:08:05+08:00"}
	//log.Fatal("Fatal message")
	//log.Panic("Panic message")

	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean") //

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!") // {"level":"warning","msg":"The group's number increased tremendously!","number":122,"omg":true,"time":"2021-09-26T16:08:05+08:00"}

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!") // {"level":"fatal","msg":"The ice breaks!","number":100,"omg":true,"time":"2021-09-26T16:08:05+08:00"}

	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other": "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")
}