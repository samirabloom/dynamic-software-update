package log

import (
	"time"
	"encoding/csv"
	"os"
	"os/signal"
	"syscall"
	"strconv"
)

// ==== PERFORMANCE - START

var entries map[int]map[string]int64 = make(map[int]map[string]int64)
var entryCounter int

func Trace(event string, startTime time.Time) {
	if entries[entryCounter] == nil {
		entries[entryCounter] = make(map[string]int64)
	}
	entries[entryCounter][event] = int64(time.Since(startTime))
}

func EndPerformanceLogEntry() {
	entryCounter++
}

var PerformanceLog = func() *csv.Writer {
	file, error := os.Create("/tmp/performance_log.csv")

	if error != nil {
		panic(error)
	}

	// New Csv writer
	writer := csv.NewWriter(file)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGKILL)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGHUP)
	signal.Notify(c, syscall.SIGQUIT)
	go func() {
		<-c
		if len(entries) > 0 {
			// header
			var (
				headers = make([]string, len(entries[0]))
				headers_counter int
			)
			for key := range entries[0] {
				headers[headers_counter] = key
				headers_counter++
			}
			returnError := writer.Write(headers)
			if returnError != nil {
				LoggerFactory().Error("Error writing headers into performance log - %s", returnError)
			}
			writer.Flush()

			// entries
			for _, entry := range entries {
				var (
					values = make([]string, len(entries))
					values_counter int
				)
				for key := range entry {
					values[values_counter] = strconv.FormatInt(entry[key], 10)
					values_counter++
				}
				returnError := writer.Write(values)
				if returnError != nil {
					LoggerFactory().Error("Error writing values into performance log - %s", returnError)
				}
			}
			writer.Flush()
		}

		file.Close()
		os.Exit(1)
	}()

	return writer
}()

// ==== PERFORMANCE - END
