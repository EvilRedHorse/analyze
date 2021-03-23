package lockcheck

import (
	"testing"
)

// TestLockcheckHelpers probes the helper functions of the lockcheck package
func TestLockcheckHelpers(t *testing.T) {
	t.Run("ContainsMutex", testContainsMutex)
	t.Run("FirstWordIs", testFirstWordIs)
	t.Run("IsManagedExported", testIsManagedExported)
	t.Run("IsMutexCall", testIsMutexCall)
	t.Run("IsStaticField", testIsStaticField)
	t.Run("IsSyncObject", testIsSyncObject)
	t.Run("ManagesOwnLocking", testManagesOwnLocking)
}

func testContainsMutex(t *testing.T) {
	t.Skip("not implemented")
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
	// Define tests
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

	// Run tests
	for _, test := range tests {
		if isManagedExported(test.name) != test.result {
			t.Error("bad", test)
		}
	}
}

func testIsMutexCall(t *testing.T) {
	t.Skip("not implemented")
}

// testIsStaticField probes the isStaticField function
func testIsStaticField(t *testing.T) {
	// Define tests
	var tests = []struct {
		name   string
		result bool
	}{
		{"field", false},
		{"staticField", true},
		{"atomicField", true},
		{"notStaticField", false},
		{"notAtomicField", false},
	}

	// Run tests
	for _, test := range tests {
		if isStaticField(test.name) != test.result {
			t.Error("bad", test)
		}
	}
}

func testIsSyncObject(t *testing.T) {
	t.Skip("not implemented")
}

// testManagesOwnLocking probes the managesOwnLocking function
func testManagesOwnLocking(t *testing.T) {
	// Define tests
	var tests = []struct {
		name   string
		result bool
	}{
		{"unexported", false},
		{"Exported", true},
		{"OtherExported", true},
		{"UnmanagedExported", false},
		{"unmanagedUnExported", false},
		{"managedMethod", true},
		{"externMethod", true},
		{"atomicMethod", false}, // Our guidelines don't talk about atomic prefix for methods
		{"threadedMethod", true},
		{"callMethod", true},
	}

	// Run tests
	for _, test := range tests {
		if managesOwnLocking(test.name) != test.result {
			t.Error("bad", test)
		}
	}
}
