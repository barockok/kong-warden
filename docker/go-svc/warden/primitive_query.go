package warden

import (
	"reflect"
	"strings"
)

const SEPARATOR = "/"
const STRING_EQUAL = "sq"
const STRING_IN = "si"
const BOOLEAN_TRUE = "bt"
const BOOLEAN_FALSE = "bf"
const EFFECT_ALLOW = "a"
const EFFECT_DENNY = "d"
const EFFECT_FORWARD_ALLOW = "fa"
const EFFECT_FORWARD_DENY = "fd"

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
	valInterface := getKeyVal(payload, key)
	if checkNilInterface(valInterface) {
		return false
	}
	if op == STRING_EQUAL {
		return stringEqual(valInterface, opval[1])
	}
	if op == STRING_IN {
		return stringInclude(valInterface, opval[1])
	}
	if op == BOOLEAN_TRUE {
		return booleanTrue(valInterface)
	}
	if op == BOOLEAN_FALSE {
		return !booleanTrue(valInterface)
	}
	return false
}

func stringEqual(value interface{}, test string) bool {
	return value.(string) == test
}
func stringInclude(value interface{}, test string) bool {
	v := value.(string)
	for _, s := range strings.Split(test, ",") {
		if v == s {
			return true
		}
	}
	return false
}
func booleanTrue(value interface{}) bool {
	return value.(bool) == true
}
func getKeyVal(payload map[string]interface{}, rawkey string) interface{} {
	keys := strings.Split(rawkey, ".")
	if len(keys) == 1 {
		return payload[keys[0]]
	}

	var val interface{} = payload[keys[0]]
	for i := 1; i < len(keys); i++ {
		val = val.(map[string]interface{})[keys[i]]
	}
	return val
}

func checkNilInterface(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}

type PrimitiveQuery struct {
	queries []string
}

func NewPrimitiveQuery(queries []string) *PrimitiveQuery {
	return &PrimitiveQuery{queries: queries}
}

func (q *PrimitiveQuery) Match(i map[string]interface{}) bool {
	for _, selector := range q.queries {
		if !evaluateSelector(i, selector) {
			return false
		}
	}
	return true
}
