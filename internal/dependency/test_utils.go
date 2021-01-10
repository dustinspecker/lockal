package dependency

import (
	"reflect"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
)

func getLogCtx() (*memory.Handler, *log.Entry) {
	log.SetLevel(log.DebugLevel)
	logHandler := memory.New()
	log.SetHandler(logHandler)

	logCtx := log.WithFields(log.Fields{
		"app": "lockal-test",
	})

	return logHandler, logCtx
}

func hasLogEntry(handler *memory.Handler, logLevel log.Level, fields log.Fields, message string) bool {
	for _, entry := range handler.Entries {
		if entry.Level == logLevel && entry.Message == message && reflect.DeepEqual(entry.Fields, fields) {
			return true
		}
	}

	return false
}
