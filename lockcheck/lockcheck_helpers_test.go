package lockcheck

import (
	"testing"
)

// TestLockcheckHelpers probes the helper functions of the lockcheck package
func TestLockcheckHelpers(t *testing.T) {

	t.Run("FirstWordIs", testFirstWordIs)
	t.Run("IsManagedExported", testIsManagedExported)
}

// testFirstWordIs probes the firstWordIs function
func testFirstWordIs(t *testing.T) {
	var tests = []struct {
		name   string
		prefix string
		result bool
	}{
		{"startsUpper", "starts", true},
		{"startsupper", "starts", false},
		{"starts", "starts", false},
	}

	for _, test := range tests {
		if firstWordIs(test.name, test.prefix) != test.result {
			t.Error("bad", test)
		}
	}
}

// testIsManagedExported probes the isManagedExported function
func testIsManagedExported(t *testing.T) {
	var tests = []struct {
		name   string
		result bool
	}{
		{"unexported", false},
		{"Exported", true},
		{"OtherExported", true},
		{"UnmanagedExported", false},
		{"unmanagedUnExported", false},
	}

	for _, test := range tests {
		if isManagedExported(test.name) != test.result {
			t.Error("bad", test)
		}
	}
}
