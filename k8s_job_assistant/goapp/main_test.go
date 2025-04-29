// main_test.go
package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	keepResources bool
	tearDown      bool
)

func TestMain(m *testing.M) {
	// Custom flag to keep resources
	flag.BoolVar(&keepResources, "keep-resources", false, "Keep test resources after test run")
	flag.BoolVar(&tearDown, "tear-down", false, "Attempt to deleteJobAndWaitForDeletion resources before testing, in case of leftovers")
	flag.Parse()

	if tearDown {
		fmt.Println("-tear-down flag set to true")
	}
	if keepResources {
		fmt.Println("-keep-resources flag set to true")
	}

	os.Exit(m.Run())
}
