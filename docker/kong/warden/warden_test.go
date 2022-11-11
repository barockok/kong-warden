package main

import (
	"fmt"
	"testing"
)

func TestMatcherTag(t *testing.T) {
	// t.Error("hello")
	type matchingCase struct {
		inputTag      string
		expectedValue string
	}

	for i, skena := range []matchingCase{
		{inputTag: "warden-action:people.view", expectedValue: "people.view"},
		{inputTag: "warden-action:transaction.view", expectedValue: "transaction.view"},
		{inputTag: "someothertag", expectedValue: ""},
	} {
		t.Run(fmt.Sprintf("skena %d , e: %s", i, skena.expectedValue), func(t *testing.T) {
			action := actionTagMatcher(skena.inputTag)
			if action != skena.expectedValue {
				t.Error("action not found")
			}
		})
	}

}

func TestActionEvaluate(t *testing.T) {
	type matchingBoolString struct {
		action     string
		rawAbility string
		result     bool
	}
	for i, skena := range []matchingBoolString{
		{
			action:     "people.view",
			rawAbility: `[{"action": "people.view"}]`,
			result:     true,
		},
	} {
		t.Run(fmt.Sprintf("skena %d, action : %s", i, skena.action), func(t *testing.T) {
			if EvaluateAbility(skena.action, skena.rawAbility) != skena.result {
				t.Error("not matching")
			}
		})
	}
}
