package iec61850

import (
	"strings"
	"testing"
)

func TestLibraryVersion(t *testing.T) {
	version := GetLibraryVersion()

	// This project requires libiec61850 v1.6.1
	// Built with R-GOOSE, R-SMV, and SNTP client enabled
	expectedVersion := "1.6.1"
	if !strings.HasPrefix(version, expectedVersion) {
		t.Fatalf("Expected version %s, got: %s. Run ./scripts/rebuild_libraries.sh to build v1.6.1", expectedVersion, version)
	}

	t.Logf("Using libiec61850 version: %s (with R-GOOSE, R-SMV, SNTP enabled)", version)
}
