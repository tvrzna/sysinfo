package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/tvrzna/go-utils/args"
)

var buildVersion string

var knownWidgets = []string{"cpu", "diskstats", "diskusage", "memory", "netspeed", "smartctl", "system", "temps", "top"}

type config struct {
	name         string
	appUrl       string
	uri          string
	port         int
	widgets      [][]string
	widgetsIndex map[string]bool
}

func loadConfig(arg []string) *config {
	c := &config{"sysinfo", "", "", 1700, make([][]string, 0), make(map[string]bool)}

	strWidgets := "cpu diskusage\\n memory system\\n temps netspeed\\n top diskstats\\n smartctl"

	args.ParseArgs(arg, func(arg, value string) {
		switch arg {
		case "-h", "--help":
			fmt.Printf(`Usage: sysinfo [options]
Options:
	-h, --help			print this help
	-v, --version			print version
	-n, --name [NAME]		name of application to be displayed
	-p, --port [PORT]		sets port for listening
	-u, --uri [URI]			sets uri (server:port) for listening
	-a, --app-url [APP_URL]		application url (if behind proxy)
	-w, --widget-layout [LAYOUT]	custom layout of widgets

Widgets:
	` + strings.Join(knownWidgets, ", ") + `

	default: '` + strWidgets + `'
`)
			os.Exit(0)
		case "-v", "--version":
			fmt.Printf("sysinfo %s\nhttps://github.com/tvrzna/sysinfo\n\nReleased under the MIT License.\n", c.getVersion())
			os.Exit(0)
		case "-n", "--name":
			c.name = value
		case "-u", "--uri":
			c.uri = value
		case "-p", "--port":
			c.port, _ = strconv.Atoi(value)
		case "-a", "--app-url":
			c.appUrl = value
		case "-w", "--widget-layout":
			strWidgets = value
		}
	})

	c.parseWidgets(strWidgets)

	return c
}

func (c *config) getServerUri() string {
	if c.uri != "" {
		return c.uri
	}
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

func (c *config) parseWidgets(cfg string) {
	cfg = strings.ReplaceAll(cfg, "\\n", "\n")
	scanner := bufio.NewScanner(strings.NewReader(cfg))

	r := regexp.MustCompile("([a-zA-Z0-1]+)")

	for scanner.Scan() {
		line := make([]string, 0)
		for _, g := range r.FindAllStringSubmatch(scanner.Text(), -1) {
			if len(g) > 0 {
				widget := strings.ToLower(g[1])
				if slices.Contains(knownWidgets, widget) {
					line = append(line, widget)
					c.widgetsIndex[widget] = true
				} else {
					log.Printf("-- unknown widget '%s' will be skipped", widget)
				}

			}
		}
		if len(line) > 0 {
			c.widgets = append(c.widgets, line)
		}
	}
}
