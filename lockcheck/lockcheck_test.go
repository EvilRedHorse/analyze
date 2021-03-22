package lockcheck_test

import (
	"testing"

	"gitlab.com/NebulousLabs/analyze/lockcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

// Test is the main test for the lockcheck package
func Test(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), lockcheck.Analyzer, "a")
}
