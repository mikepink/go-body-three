package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
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

func getWsServer() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling new WebSocket request. Upgrading request...")
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			panic(err)
		}

		var wg sync.WaitGroup
		positionChan := make(chan *Frame, 60)
		endSimulationChan := make(chan bool, 1)
		wg.Add(1)
		go Simulator(wg, positionChan, endSimulationChan)
		wg.Add(1)
		go func(positionChan <-chan *Frame) {
			defer wg.Done()
			defer conn.Close()
			fmt.Println("Initializing simulation engine")

			fmt.Println("Awaiting client data...")
			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					log.Println("wsutil.ReadClientData: ", err)
					os.Exit(0)
					return
				}
				fmt.Printf("Received data[%d] from client: %s\n", op, msg)
				bufSize := 60
				positionBuffer := make([]*Frame, bufSize)
				for i := 0; i < bufSize; i++ {
					select {
					case positionChanMsg := <-positionChan:
						positionBuffer[i] = positionChanMsg
					case endSimulationChanMsg := <-endSimulationChan:
						fmt.Printf("End Simulation Message: %t\n", endSimulationChanMsg)
						return
					}
				}
				jsonBytes, jsonerr := json.Marshal(positionBuffer)
				if jsonerr != nil {
					log.Println("json.Marshall: ", jsonerr)
				}
				responseBuffer := jsonBytes
				err = wsutil.WriteServerMessage(conn, ws.OpText, responseBuffer)
				if err != nil {
					log.Println("wsutil.WriteServerMessage: ", err)
					return
				}
				fmt.Println("Responsed to client...")
			}
		}(positionChan)

		wg.Wait()
	}
}

func init_server() {
	fmt.Println("Generating static resources map...")
	makeStaticResoucesMap(staticResoucesMap)
	fmt.Println("Initializing HTTP server...")
	srv := &http.Server{Addr: ":8822"}
	http.HandleFunc("/", MainWebServer)
	http.HandleFunc("/sim", getWsServer())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := srv.ListenAndServeTLS("server.crt", "server.key")
		if err != nil {
			log.Fatal("ListenAndServerTLS: ", err)
		}
		wg.Done()
	}()
	wg.Wait()
}

func main() {
	init_server()
}
