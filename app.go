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

	for webPath, localPath := range srMap {
		srMap[webPath] = filepath.Join(rootPath, STATIC_ROOT, localPath)
	}
}

func getContentTypeForFile(filePath string) string {
	if strings.HasSuffix(filePath, ".html") {
		return "text/html"
	} else if strings.HasSuffix(filePath, ".css") {
		return "text/css"
	} else if strings.HasSuffix(filePath, ".js") {
		return "text/javascript"
	} else {
		return "text/plain"
	}
}

func MainWebServer(w http.ResponseWriter, req *http.Request) {
	staticResourcePath, foundStaticResource := staticResoucesMap[req.URL.Path]
	if foundStaticResource {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", getContentTypeForFile(staticResourcePath))
		fileBytes, err := ioutil.ReadFile(staticResoucesMap[req.URL.Path])
		if err != nil {
			panic(err)
		}
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
