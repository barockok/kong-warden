package main

import "strings"

const SEPARATOR = "/"
const STRING_EQUAL = "sq"
const STRING_IN = "si"
const BOOLEAN_TRUE = "bt"
const BOOLEAN_FALSE = "bf"
const EFFECT_ALLOW = "a"
const EFFECT_DENNY = "d"

func AttributeMatch(payload map[string]interface{}, resourceMatcher []string) bool {
	for _, selector := range resourceMatcher {
		if !evaluateSelector(payload, selector) {
			return false
		}
	}
	return true
}

func evaluateSelector(payload map[string]interface{}, pattern string) bool {
	segments := strings.Split(pattern, SEPARATOR)
	effect := segments[0]
	matcher := segments[1:]

	for i, attrpair := range matcher {
		if i%2 == 0 {
			if EFFECT_DENNY == effect && evaluateAttribute(payload, attrpair, matcher[i+1]) {
				return false
			}
			if EFFECT_ALLOW == effect && !evaluateAttribute(payload, attrpair, matcher[i+1]) {
				return false
			}
		}
	}
	return true
}

func evaluateAttribute(payload map[string]interface{}, key, rawVal string) bool {
	opval := strings.Split(rawVal, ":")

	op := opval[0]

	if op == STRING_EQUAL {
		return stringEqual(getKeyVal(payload, key).(string), opval[1])
	}
	if op == STRING_IN {
		return stringInclude(getKeyVal(payload, key).(string), opval[1])
	}
	if op == BOOLEAN_TRUE {
		return booleanTrue(getKeyVal(payload, key).(bool))
	}
	if op == BOOLEAN_FALSE {
		return !booleanTrue(getKeyVal(payload, key).(bool))
	}
	return false
}

func stringEqual(value, test string) bool {
	return value == test
}
func stringInclude(value, test string) bool {
	for _, s := range strings.Split(test, ",") {
		if value == s {
			return true
		}
	}
	return false
}
func booleanTrue(value bool) bool {
	return value == true
}
func getKeyVal(payload map[string]interface{}, key string) interface{} {
	return payload[key]
}
