package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
)

//go:embed www
var www embed.FS

type context struct {
	conf      *config
	webServer *http.Server
	interrupt chan bool
	layout    *template.Template
}

type PageContext struct {
	c       *context
	Name    string
	Version string
}

func (p *PageContext) UrlFor(path string) string {
	if p.c.conf.getAppUrl() == "" {
		return path
	}
	return p.c.conf.getAppUrl() + "/" + path
}

func newContext(conf *config) *context {
	result := &context{conf: conf, interrupt: make(chan bool)}

	tpl, err := template.ParseFS(www, "www/template.html")
	if err != nil {
		log.Fatal(err)
	}
	result.layout = tpl

	return result
}

func (c *context) runWebServer() {
	mux := http.NewServeMux()

	serverRoot, _ := fs.Sub(www, "www")
	mux.Handle("/static/", http.FileServer(http.FS(serverRoot)))
	mux.HandleFunc("/sysinfo.json", HandleSysinfoData)
	mux.HandleFunc("/", c.handleIndex)

	c.webServer = &http.Server{Handler: mux, Addr: c.conf.getServerUri()}

	log.Print("-- sysinfo started on ", c.conf.getServerUri())
	if err := c.webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Print("-- sysinfo start failed: ", err)
		c.interrupt <- true
	} else {
		log.Print("-- sysinfo finished")
	}
}

func (c *context) handleStop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	select {
	case <-ch:
	case <-c.interrupt:
	}

	log.Print("-- stopping sysinfo")
	if c.webServer != nil {
		c.webServer.Close()
	}
}

func (c *context) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" || r.URL.Path == "index.html" {
		p := &PageContext{c: c, Name: c.conf.name, Version: c.conf.getVersion()}
		w.Header().Set("content-type", "text/html")

		if err := c.layout.Execute(w, p); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
