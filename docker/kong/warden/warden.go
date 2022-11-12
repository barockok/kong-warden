package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"golang.org/x/exp/slices"
)

const actionTagRE = "warden-action:(.*)"
const errorActionTagNotFound = "Action Tag Not Found"
const errorAbilityNotFound = "Ability Not Found"
const HeaderWardenAbilities = "X-Warden-Abilities"

func actionTagMatcher(s string) string {
	re := regexp.MustCompile(actionTagRE)
	m := re.FindStringSubmatch(s)
	if len(m) < 1 {
		return ""
	}
	return m[1]
}

func main() {
	server.StartServer(New, Version, Priority)
}

var Version = "0.2"
var Priority = 1

type Config struct {
	ActionName string `json:"actionName"`
}

func New() interface{} {
	return &Config{}
}

type WardenAbility struct {
	Action   string   `json:"a"`
	Selector []string `json:"s"`
}

func FindAbility(action string, abilities []WardenAbility) (WardenAbility, error) {
	for _, ability := range abilities {
		if action == ability.Action {
			return ability, nil
		}
	}
	return WardenAbility{}, errors.New(errorAbilityNotFound)
}
func EvaluateSelector(kong *pdk.PDK, ability WardenAbility) bool {
	httpMethod, err := kong.Request.GetMethod()
	if err != nil {
		return true
	}

	contentType, err := kong.Request.GetHeader("Content-Type")
	if err != nil {
		return true
	}

	if slices.Index([]string{"post", "put", "patch"}, strings.ToLower(httpMethod)) > -1 && contentType == "application/json" {
		rawBody, err := kong.Request.GetRawBody()
		if err != nil {
			return true
		}
		var payload map[string]interface{}
		json.Unmarshal(rawBody, &payload)
		return AttributeMatch(payload, ability.Selector)
	}

	return true
}

func GetRouteAction(kong *pdk.PDK) (string, error) {
	route, err := kong.Router.GetRoute()
	if err != nil {
		return "", err
	}
	for _, oriTag := range route.Tags {
		tag := actionTagMatcher(oriTag)
		if tag != "" {
			return tag, nil
		}
	}
	return "", errors.New(errorActionTagNotFound)
}

func (conf *Config) Access(kong *pdk.PDK) {
	action, err := GetRouteAction(kong)
	if err != nil {
		kong.Log.Err(fmt.Sprintf("[warden] action not found in route tag %v", err))
		kong.Response.ExitStatus(403)
		return
	}

	rawWardenAbility, err := kong.Request.GetHeader("X-Warden-Abilities")
	if err != nil {
		kong.Log.Err(fmt.Sprintf("[warden] ability header not present %v", err))
		kong.Response.ExitStatus(403)
		return
	}

	kong.Response.SetHeader("x-warden-action", action)
	kong.Response.SetHeader("x-warden-ability", rawWardenAbility)

	var parsedWardenAbility []WardenAbility
	json.Unmarshal([]byte(rawWardenAbility), &parsedWardenAbility)

	ability, err := FindAbility(action, parsedWardenAbility)
	if err != nil {
		kong.Response.ExitStatus(403)
		return
	}

	if !EvaluateSelector(kong, ability) {
		kong.Response.ExitStatus(403)
		return
	}
}
