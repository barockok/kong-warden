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
const errorPermissionNotFound = "Permission Not Found"
const HeaderWardenPermissions = "X-Warden-Permissions"
const HeaderWardenPermissionsForward = "X-Warden-Permissions-Forward"

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

type WardenPermission struct {
	Action   string   `json:"a"`
	Selector []string `json:"s"`
}

func FindPermission(action string, abilities []WardenPermission) (WardenPermission, error) {
	for _, permission := range abilities {
		if action == permission.Action {
			return permission, nil
		}
	}
	return WardenPermission{}, errors.New(errorPermissionNotFound)
}

func ForwadEffect(kong *pdk.PDK, permission WardenPermission) {
	headerVal := []string{}
	for _, s := range permission.Selector {
		if strings.HasPrefix(s, EFFECT_FORWARD_ALLOW) {
			headerVal = append(headerVal, fmt.Sprintf("%s%s", EFFECT_ALLOW, strings.TrimPrefix(s, EFFECT_FORWARD_ALLOW)))
		}
		if strings.HasPrefix(s, EFFECT_FORWARD_DENY) {
			headerVal = append(headerVal, fmt.Sprintf("%s%s", EFFECT_DENNY, strings.TrimPrefix(s, EFFECT_FORWARD_DENY)))
		}
	}
	headerValStr, err := json.Marshal(headerVal)
	if err != nil {
		kong.Log.Warn(fmt.Sprintf("[warden] cannot decode forwardEffect %v", err))
		return
	}
	kong.ServiceRequest.AddHeader(HeaderWardenPermissionsForward, string(headerValStr))
}
func EvaluateSelector(kong *pdk.PDK, permission WardenPermission) bool {
	httpMethod, err := kong.Request.GetMethod()
	if err != nil {
		return true
	}

	ForwadEffect(kong, permission)
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
		return AttributeMatch(payload, permission.Selector)
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

	rawWardenPermission, err := kong.Request.GetHeader("X-Warden-Permissions")
	if err != nil {
		kong.Log.Err(fmt.Sprintf("[warden] permission header not present %v", err))
		kong.Response.ExitStatus(403)
		return
	}

	kong.Response.SetHeader("x-warden-action", action)
	kong.Response.SetHeader("x-warden-permission", rawWardenPermission)

	var parsedWardenPermission []WardenPermission
	json.Unmarshal([]byte(rawWardenPermission), &parsedWardenPermission)

	permission, err := FindPermission(action, parsedWardenPermission)
	if err != nil {
		kong.Response.ExitStatus(403)
		return
	}

	if !EvaluateSelector(kong, permission) {
		kong.Response.ExitStatus(403)
		return
	}
}
