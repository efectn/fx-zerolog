package fxzerolog

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"text/tabwriter"
)

// From: https://github.com/gofiber/fiber/blob/master/utils/assertions.go
func assertEqual(tb testing.TB, expected, actual interface{}, description ...string) {
	if tb != nil {
		tb.Helper()
	}

	if reflect.DeepEqual(expected, actual) {
		return
	}

	aType := "<nil>"
	bType := "<nil>"

	if expected != nil {
		aType = reflect.TypeOf(expected).String()
	}
	if actual != nil {
		bType = reflect.TypeOf(actual).String()
	}

	testName := "AssertEqual"
	if tb != nil {
		testName = tb.Name()
	}

	_, file, line, _ := runtime.Caller(1)

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 5, ' ', 0)
	fmt.Fprintf(w, "\nTest:\t%s", testName)
	fmt.Fprintf(w, "\nTrace:\t%s:%d", filepath.Base(file), line)
	if len(description) > 0 {
		fmt.Fprintf(w, "\nDescription:\t%s", description[0])
	}
	fmt.Fprintf(w, "\nExpect:\t%v\t(%s)", expected, aType)
	fmt.Fprintf(w, "\nResult:\t%v\t(%s)", actual, bType)

	result := ""
	if err := w.Flush(); err != nil {
		result = err.Error()
	} else {
		result = buf.String()
	}

	if tb != nil {
		tb.Fatal(result)
	} else {
		log.Fatal(result)
	}
}
