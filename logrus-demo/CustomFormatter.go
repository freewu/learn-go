package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type CustomFormatter struct {
}

// 不自动添加 level / msg / time
func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	// Note this doesn't include Time, Level and Message which are available on
	// the Entry. Consult `godoc` on information about those fields or read the
	// source of the official loggers.
	fmt.Println(entry.Level)
	fmt.Println(entry.Time)
	fmt.Println(entry.Message)
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %w", err)
	}
	return append(serialized, '\n'), nil
}

func main() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&CustomFormatter{})

	log.Trace("Trace message") // {}
	log.Debug("Debug message") // {}
	log.Info("Info message") // {}
	log.Warn("Warn message") // {}
	log.Error("Error message") // {}
}