package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const STATIC_ROOT = "/web/static"

var staticResoucesMap = make(map[string]string)

func makeStaticResoucesMap(srMap map[string]string) {
	rootPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(rootPath)
	srMap["/"] = "index.html"
	srMap["/favicon.ico"] = "favicon.ico"
	srMap["/js/app.js"] = "js/app.js"
	srMap["/js/three.js"] = "js/three.js"
	srMap["/styles/app.css"] = "styles/app.css"

	for webPath, localPath := range srMap {
		srMap[webPath] = filepath.Join(rootPath, STATIC_ROOT, localPath)
	}
}

func getContentTypeForFile(filePath string) string {
	if strings.HasSuffix(filePath, ".html") {
		return "text/html; charset=utf-8"
	} else if strings.HasSuffix(filePath, ".css") {
		return "text/css; charset=utf-8"
	} else if strings.HasSuffix(filePath, ".ico") {
		return "image/x-icon"
	} else if strings.HasSuffix(filePath, ".js") {
		return "text/javascript; charset=utf-8"
	} else {
		return "text/plain; charset=utf-8"
	}
}

func MainWebServer(w http.ResponseWriter, req *http.Request) {
	staticResourcePath, foundStaticResource := staticResoucesMap[req.URL.Path]
	if foundStaticResource {
		fileBytes, err := ioutil.ReadFile(staticResoucesMap[req.URL.Path])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			panic(err)
		}

		w.Header().Set("Content-Type", getContentTypeForFile(staticResourcePath))
		w.WriteHeader(http.StatusOK)

		w.Write(fileBytes)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func init_server() {
	fmt.Println("Generating static resources map...")
	makeStaticResoucesMap(staticResoucesMap)
	fmt.Println("Initializing HTTP server...")
	http.HandleFunc("/", MainWebServer)
	err := http.ListenAndServeTLS(":8822", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServerTLS: ", err)
	}
}

func main() {
	init_server()
}
