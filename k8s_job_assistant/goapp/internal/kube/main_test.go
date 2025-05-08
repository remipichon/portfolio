// main_test.go
package kube

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	TearDown      bool
	KeepResources bool
	Kubeconfig    string
)

func TestMain(m *testing.M) {
	// Custom flag to keep resources
	flag.BoolVar(&KeepResources, "keep-resources", false, "Keep test resources after test run")
	flag.BoolVar(&TearDown, "tear-down", false, "Attempt to deleteJobAndWaitForDeletion resources before testing, in case of leftovers")
	flag.StringVar(&Kubeconfig, "kubeconfig", "", "path to kubeconfig for test")
	flag.Parse()

	if TearDown {
		fmt.Println("-tear-down flag set to true")
	}
	if KeepResources {
		fmt.Println("-keep-resources flag set to true")
	}

	os.Exit(m.Run())
}
