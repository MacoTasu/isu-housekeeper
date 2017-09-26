package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"

	webhooks "gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/github"
	yaml "gopkg.in/yaml.v2"
)

const (
	path = "/webhooks"
	port = 3016
)

func main() {

	// webhook secrets
	buf, err := ioutil.ReadFile(".secrets.yaml")
	if err != nil {
		panic(err)
	}

	var config github.Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		panic(err)
	}

	hook := github.New(&config)
	hook.RegisterEvents(HandleMultiple, github.PushEvent)

	err = webhooks.Run(hook, ":"+strconv.Itoa(port), path)
	if err != nil {
		fmt.Println(err)
	}
}

// HandleMultiple handles multiple GitHub events
func HandleMultiple(payload interface{}, header webhooks.Header) {

	switch payload.(type) {

	case github.PushPayload:
		push := payload.(github.PushPayload)
		fmt.Printf("%+v", push)

		assined := regexp.MustCompile("refs/heads/(.*)")
		group := assined.FindStringSubmatch(push.Ref)
		if len(group) >= 2 {
			err := exec.Command("./scripts/deploy.sh", group[1]).Run()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
