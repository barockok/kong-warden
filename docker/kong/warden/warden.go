package main

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

const actionTagRE = "warden-action:(.*)"
const errorActionTagNotFound = "Action Tag Not Found"
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
	Message string
}

func New() interface{} {
	return &Config{}
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

type WardenAbility struct {
	Action string `json:"action"`
}

func EvaluateAbility(action, rawWardenAbility string) bool {
	var abilities []WardenAbility
	json.Unmarshal([]byte(rawWardenAbility), &abilities)
	for _, ability := range abilities {
		if action == ability.Action {
			return true
		}
	}
	return false
}

func (conf Config) Access(kong *pdk.PDK) {
	action, err := GetRouteAction(kong)
	if err != nil {
		log.Printf("[Warning] Get Route Action, %s", err.Error())
	}

	warWardenAbility, err := kong.Request.GetHeader("X-Warden-Abilities")

	if err != nil {
		log.Printf("[Warning] Get Warnden Ability, %s", err.Error())
	}

	kong.Response.SetHeader("x-warden-action", action)
	kong.Response.SetHeader("x-warden-ability", warWardenAbility)

	if !EvaluateAbility(action, warWardenAbility) {
		kong.Response.ExitStatus(403)
	}
}

func (conf Config) Response(kong *pdk.PDK) {
	action, err := GetRouteAction(kong)
	if err != nil {
		log.Printf("[Warning] Get Route Action, %s", err.Error())
	}

	warWardenAbility, err := kong.Request.GetHeader("X-Warden-Abilities")

	if err != nil {
		log.Printf("[Warning] Get Warnden Ability, %s", err.Error())
	}

	kong.Response.SetHeader("x-warden-action", action)
	kong.Response.SetHeader("x-warden-ability", warWardenAbility)

	if !EvaluateAbility(action, warWardenAbility) {
		kong.Response.ExitStatus(403)
	}
}
