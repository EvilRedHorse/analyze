package lockcheck

import (
	"go/types"

	"testing"
)

// TestLockcheckHelpers probes the helper functions of the lockcheck package
func TestLockcheckHelpers(t *testing.T) {
	t.Run("ContainsMutex", testContainsMutex)
	t.Run("FirstWordIs", testFirstWordIs)
	t.Run("IsManagedExported", testIsManagedExported)
	t.Run("IsMutexCall", testIsMutexCall)
	t.Run("IsMutexType", testIsMutexType)
	t.Run("IsStaticField", testIsStaticField)
	t.Run("IsSyncObject", testIsSyncObject)
	t.Run("ManagesOwnLocking", testManagesOwnLocking)
}

func testContainsMutex(t *testing.T) {
	// Need to figure out how to satisfy the types.Object interface
	t.Skip("not implemented")
}

// testFirstWordIs probes the firstWordIs function
func testFirstWordIs(t *testing.T) {
	var tests = []struct {
		name   string
		prefix string
		result bool
	}{
		// Valid cases
		{"startsUpper", "starts", true},

		// Invalid cases
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
		// Valid cases
		{"Exported", true},
		{"OtherExported", true},

		// Invalid cases
		{"unexported", false},
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
	// Need to figure out how to satisfy the types.Object interface
	t.Skip("not implemented")
}

// testIsMutexType probes the isMutexType function
func testIsMutexType(t *testing.T) {
	var foo fooType
	if isMutexType(foo.Underlying()) {
		t.Error("bad")
	}

	var fooMu fooMutex
	if !isMutexType(fooMu.Underlying()) {
		t.Error("bad")
	}

	var fooMuNot fooMutexNot
	if isMutexType(fooMuNot.Underlying()) {
		t.Error("bad")
	}
}

// testIsStaticField probes the isStaticField function
func testIsStaticField(t *testing.T) {
	// Define tests
	var tests = []struct {
		name   string
		result bool
	}{
		// Valid cases
		{"staticField", true},
		{"atomicField", true},

		// Invalid cases
		{"externField", false},
		{"field", false},
		{"notStaticField", false},
		{"notAtomicField", false},
		{"noExternField", false},
	}

	// Run tests
	for _, test := range tests {
		if isStaticField(test.name) != test.result {
			t.Error("bad", test)
		}
	}
}

// testIsSyncObject probes the isSyncObject function
func testIsSyncObject(t *testing.T) {
	// Define variable types to test
	var mu syncMutex
	var rwMu syncRWMutex
	var wg syncWaitGroup
	var siaRW siaSyncRWMutex
	var siaTry siaSyncTryMutex
	var siaTryRW siaSyncTryRWMutex
	var tg threadgroup

	var foo fooType
	var fooMu fooMutex
	var fooMuNot fooMutexNot

	// Define tests
	var tests = []struct {
		object string
		t      types.Type
		result bool
	}{
		// Valid cases
		{mu.String(), mu.Underlying(), true},
		{rwMu.String(), rwMu.Underlying(), true},
		{wg.String(), wg.Underlying(), true},
		{siaRW.String(), siaRW.Underlying(), true},
		{siaTry.String(), siaTry.Underlying(), true},
		{siaTryRW.String(), siaTryRW.Underlying(), true},
		{tg.String(), tg.Underlying(), true},

		// Invalid cases
		{foo.String(), foo.Underlying(), false},
		{fooMu.String(), fooMu.Underlying(), false},
		{fooMuNot.String(), fooMuNot.Underlying(), false},
	}

	// Run Tests
	for _, test := range tests {
		if isSyncObject(test.t) != test.result {
			t.Error("bad", test.object)
		}
	}
}

// testManagesOwnLocking probes the managesOwnLocking function
func testManagesOwnLocking(t *testing.T) {
	// Define tests
	var tests = []struct {
		name   string
		result bool
	}{
		// Valid cases
		{"Exported", true},
		{"OtherExported", true},
		{"managedMethod", true},
		{"externMethod", true},
		{"threadedMethod", true},
		{"callMethod", true},

		// Invalid cases
		{"unexported", false},
		{"UnmanagedExported", false},
		{"unmanagedUnExported", false},
		{"atomicMethod", false}, // Our guidelines don't talk about atomic prefix for methods
	}

	// Run tests
	for _, test := range tests {
		if managesOwnLocking(test.name) != test.result {
			t.Error("bad", test)
		}
	}
}
