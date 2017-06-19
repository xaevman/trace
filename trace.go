package trace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/xaevman/log"
)

const TraceDir = "trace"

var (
	DebugLogger log.DebugLogger
	ErrorLogger log.ErrorLogger
)

func Log(traceName string, context interface{}) {
	ctxJson, err := json.MarshalIndent(context, "", "    ")
	if err != nil {
		_err("Error generating context JSON: %v", err)
	}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("**** Start Trace::%s ****\n", traceName))
	buffer.WriteString("::Stack Trace::\n")
	buffer.Write(debug.Stack())
	buffer.WriteString("\n")

	if err == nil {
		buffer.WriteString("::Context::\n")
		buffer.Write(ctxJson)
		buffer.WriteString("\n")
	}

	buffer.WriteString(fmt.Sprintf("**** End Trace::%s ****\n", traceName))

	traceStr := buffer.String()

	traceFile := fmt.Sprintf("%s.%s.log", traceName, _fmtTime())
	tracePath := filepath.Join(TraceDir, traceFile)

	err = ioutil.WriteFile(tracePath, []byte(traceStr), 0660)
	if err != nil {
		_err("Error writing trace file %s: %v", tracePath, err)
	}
}

func init() {
	err := os.MkdirAll(TraceDir, 0660)
	if err != nil {
		panic(err)
	}
}

func _fmtTime() string {
	return time.Now().Format("20060102.150405.999999999")
}

func _err(fmt string, v ...interface{}) {
	if ErrorLogger != nil {
		ErrorLogger.Error(fmt, v...)
	}
}

func _print(fmt string, v ...interface{}) {
	if DebugLogger != nil {
		DebugLogger.Debug(fmt, v...)
	}
}
