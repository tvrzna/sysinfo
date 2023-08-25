package main

import (
	"os"
)

func main() {
	c := newContext(loadConfig(os.Args))
	go c.runWebServer()
	c.handleStop()
}
