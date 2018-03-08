// Package testassist holds testing helpers
package testassist

import "testing"

// TestCase represents a test case where there is some input and expected output.
type TestCase struct {
	Expected interface{}
	Input    interface{}
}

// Argument has some name and a value, it is provided methods during invocation
type Argument struct {
	Name  string
	Value interface{}
}

// VerifyEqual checks whenever the expected argument and the actual one are equal.
// If they are not it fails the testing with the provided message
func VerifyEqual(expected, actual, message string, t *testing.T) {
	if actual != expected {
		t.Fatal(message)
	}
}

// VerifyCondition checks if condition is satisfied.
// If it is not it fails the testing with the provided message
func VerifyCondition(condition bool, message string, t *testing.T) {
	if !condition {
		t.Fatal(message)
	}
}
