package commons

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/smartrecruiters/docker-bakery/bakery/commons/testassist"
)

func TestRemoveLast(t *testing.T) {
	testCases := []testassist.TestCase{
		{Expected: []string{}, Input: []string{"a"}},
		{Expected: []string{"a"}, Input: []string{"a", "b"}},
		{Expected: []string{"a", "b"}, Input: []string{"a", "b", "c"}},
	}

	for i, tc := range testCases {
		actual := RemoveLast(tc.Input.([]string))
		testassist.VerifyCondition(reflect.DeepEqual(tc.Expected.([]string), actual), fmt.Sprintf("TestCase: %d Expected %s, got %s", i, tc.Expected, actual), t)
	}
}

func TestRemoveIndex(t *testing.T) {
	testCases := []testassist.TestCase{
		{Expected: []string{"b", "c"}, Input: []testassist.Argument{{Name: "slice", Value: []string{"a", "b", "c"}}, {Name: "index", Value: 0}}},
		{Expected: []string{"a", "c"}, Input: []testassist.Argument{{Name: "slice", Value: []string{"a", "b", "c"}}, {Name: "index", Value: 1}}},
		{Expected: []string{"a", "b"}, Input: []testassist.Argument{{Name: "slice", Value: []string{"a", "b", "c"}}, {Name: "index", Value: 2}}},
	}

	for i, tc := range testCases {
		args := tc.Input.([]testassist.Argument)
		actual := RemoveIndex(args[0].Value.([]string), args[1].Value.(int))
		testassist.VerifyCondition(reflect.DeepEqual(tc.Expected.([]string), actual), fmt.Sprintf("TestCase: %d Expected %s, got %s", i, tc.Expected, actual), t)
	}
}
