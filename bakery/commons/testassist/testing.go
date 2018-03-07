package testassist

import "testing"

type TestCase struct {
	Expected interface{}
	Input    interface{}
}

type Argument struct {
	Name  string
	Value interface{}
}

func VerifyEqual(expected, actual, message string, t *testing.T) {
	if actual != expected {
		t.Fatal(message)
	}
}

func VerifyCondition(condition bool, message string, t *testing.T) {
	if !condition {
		t.Fatal(message)
	}
}
