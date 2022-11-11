package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

const actionTagRE = "warden-action:(.*)"
const errorActionTagNotFound = "Action Tag Not Found"

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

func (conf Config) Access(kong *pdk.PDK) {
	host, err := kong.Request.GetHeader("host")
	if err != nil {
		log.Printf("Error reading 'host' header: %s", err.Error())
	}
	action, err := GetRouteAction(kong)
	if err != nil {
		log.Printf("[Warning] Get Route Action, %s", err.Error())
	}

	message := conf.Message
	if message == "" {
		message = "hello"
	}
	kong.Response.SetHeader("x-hello-from-go", fmt.Sprintf("Go says %s to %s", message, host))
	kong.Response.SetHeader("x-warden-action", action)
}
