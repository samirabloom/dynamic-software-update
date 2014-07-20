package proxy_c

import (
	"time"
	"encoding/csv"
	"os"
	"os/signal"
	"syscall"
	"strconv"
)

// ==== PERFORMANCE - START

type performance struct {
read     *int64
route    *int64
write    *int64
complete *int64
}

func trace(startTime time.Time, result *int64) {
	*result = int64(time.Since(startTime))
}

var performanceLog = func() *csv.Writer {
	file, error := os.Create("performance_log.csv")

	if error != nil {
		panic(error)
	}

	// New Csv writer
	writer := csv.NewWriter(file)

	// Headers
	var new_headers = []string{"count", "read", "route", "write", "complete"}
	returnError := writer.Write(new_headers)
	if returnError != nil {
		loggerFactory().Error("Error writing headers into performance log - %s", returnError)
	}
	writer.Flush()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGKILL)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGHUP)
	signal.Notify(c, syscall.SIGQUIT)
	go func() {
		<-c
		file.Close()
		os.Exit(1)
	}()

	return writer
}()

func writePerformanceLogEntry(context *chunkContext) {
	performanceLog.Write([]string{
	strconv.FormatInt(*context.performance.read, 10),
	strconv.FormatInt(*context.performance.route, 10),
	strconv.FormatInt(*context.performance.write, 10),
	strconv.FormatInt(*context.performance.complete, 10)})
	performanceLog.Flush()
}

// ==== PERFORMANCE - END
