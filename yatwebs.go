/*
 * yatwebs - Yet Another Tiny Web Server (V 0.9.0)
 *
 * Just for learning to play with net/http
 * https://golang.org/doc/articles/wiki/
 * https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/03.0.html
 *
 * (c) 2015 by Bj  rn Winkler
 *
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// this http server never should add "StatusNotModified" to the http header
// so we don't use http.ServeFile()
func serveFile(w http.ResponseWriter, r *http.Request, name string) {
	content := ""

	nameFileInfo, err := os.Stat("." + name)

	if err != nil {
		http.Error(w, name+" not found", 404)
		return
	}

	if nameFileInfo.IsDir() {
		if name[len(name)-1] != '/' {
			http.Redirect(w, r, name+"/", 301)
			return
		}

		fileInfos, err := ioutil.ReadDir("." + name)

		if err != nil {
			http.Error(w, "dir "+name+" not found", 404)
			return
		}

		// todo: better string concatination? => Reader?
		content = ""
		content = content + "<html><body><pre>\n"

		for _, fileInfo := range fileInfos {
			thisName := fileInfo.Name()

			if fileInfo.IsDir() {
				thisName = thisName + "/"
			}

			content = content + "<a href=\"" + thisName + "\">" + thisName + "</a>\n"
		}

		content = content + "</pre></body></html>"

	} else {
		bytes, err := ioutil.ReadFile("." + name)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		content = string(bytes)
	}

	http.ServeContent(w, r, name, time.Now(), strings.NewReader(content))
}

func fileServer(w http.ResponseWriter, r *http.Request) {
	var path = r.URL.Path

	// print log
	fmt.Println("[" + time.Now().Format("2006-01-02 15:04:05") + "] " + path)

	// serve file
	//http.ServeFile(w, r, path)
	serveFile(w, r, path)
}

func main() {
	port := ""

	if len(os.Args) > 1 {
		port = ":" + os.Args[1]
	}

	fmt.Println("yatwebs - Yet Another Tiny Web Server (version 0.9.1)")
	fmt.Println()
	fmt.Println("starting http server on \"http://localhost" + port + "/\".")
	fmt.Println("press CTRL+c for stopping.")
	fmt.Println()

	http.HandleFunc("/", fileServer)
	http.ListenAndServe(port, nil)
}
