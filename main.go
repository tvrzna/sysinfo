package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

var buildVersion string

//go:embed www
var www embed.FS

func main() {
	serverRoot, err := fs.Sub(www, "www")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(serverRoot)))
	http.HandleFunc("/sysinfo.json", HandleSysinfoData)
	http.ListenAndServe(":1700", nil)
}
