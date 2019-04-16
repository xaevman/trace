//  ---------------------------------------------------------------------------
//
//  all_test.go
//
//  Copyright (c) 2014, Jared Chavez.
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

package trace

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"

    "github.com/xaevman/app"
)

var (
    testLogger TestLogger
    errChan    = make(chan bool, 1)
)

type TestTraceData struct {
    Title     string
    Timestamp time.Time
    TestValue int64
}

type TestLogger struct{}

func (tl TestLogger) Debug(format string, v ...interface{}) {
    fmt.Printf(format, v...)
}

func (tl TestLogger) Error(format string, v ...interface{}) {
    fmt.Printf(format, v...)
    close(errChan)
}

func TestTrace(t *testing.T) {
    ErrorLogger = testLogger
    DebugLogger = testLogger

    traceDir := filepath.Join(app.GetExeDir(), TraceDir)
    _, err := os.Stat(traceDir)
    if err != nil && os.IsNotExist(err) {
        t.Errorf("Trace directory not initialized")
    }

    traceData := &TestTraceData{
        Title:     "This is my test data",
        Timestamp: time.Now(),
        TestValue: 123456789,
    }

    jsData, err := json.MarshalIndent(traceData, "", "    ")
    if err != nil {
        t.Error(err)
    }

    tracePath := Log(
        "MyTestTrace",
        traceData,
    )

    select {
    case <-errChan:
        t.Errorf("Error logging trace data")
    default:
    }

    fileData, err := ioutil.ReadFile(tracePath)
    if err != nil {
        t.Error(err)
    }

    if !strings.Contains(string(fileData), string(jsData)) {
        t.Errorf("Trace json data not found in trace file")
    }
}
