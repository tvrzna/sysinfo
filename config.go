package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tvrzna/go-utils/args"
)

var buildVersion string

type config struct {
	name   string
	appUrl string
	port   int
}

func loadConfig(arg []string) *config {
	c := &config{"sysinfo", "", 1700}

	args.ParseArgs(arg, func(arg, value string) {
		switch arg {
		case "-h", "--help":
			fmt.Printf("Usage: sysinfo [options]\nOptions:\n\t-h, --help\t\t\tprint this help\n\t-v, --version\t\t\tprint version\n\t-n, --name [NAME]\t\tname of application to be displayed\n\t-p, --port [PORT]\t\tsets port for listening\n\t-a, --app-url [APP_URL]\t\tapplication url (if behind proxy)\n")
			os.Exit(0)
		case "-v", "--version":
			fmt.Printf("sysinfo %s\nhttps://github.com/tvrzna/sysinfo\n\nReleased under the MIT License.\n", c.getVersion())
			os.Exit(0)
		case "-n", "--name":
			c.name = value
		case "-p", "--port":
			c.port, _ = strconv.Atoi(value)
		case "-a", "--app-url":
			c.appUrl = value
		}
	})

	return c
}

func (c *config) getServerUri() string {
	return "127.0.0.1:" + strconv.Itoa(c.port)
}

func (c *config) getAppUrl() string {
	return c.appUrl
}

func (c *config) getVersion() string {
	if buildVersion == "" {
		return "develop"
	}
	return buildVersion
}
