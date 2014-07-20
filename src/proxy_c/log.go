package proxy_c

import (
	"log"
	"os"
	logging "github.com/op/go-logging"
)

// ==== LOGGER - START

var logLevel *string

var loggerFactory = func() func() *logging.Logger {
	var logg *logging.Logger = nil

	return func() *logging.Logger {
		if logg == nil {
			logg = logging.MustGetLogger("main")

			// Customize the output format
			logging.SetFormatter(logging.MustStringFormatter("%{level:8s} - %{message}"))

			// Setup one stdout and one syslog backend
			logBackend := logging.NewLogBackend(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
			logBackend.Color = true

			// Combine them both into one logging backend
			logging.SetBackend(logBackend)

			// set log level
			level, _ := logging.LogLevel("WARN")
			if logLevel != nil {
				level, _ = logging.LogLevel(*logLevel)
			}
			logging.SetLevel(level, "main")
		}
		return logg
	}
}()

// ==== LOGGER - END
