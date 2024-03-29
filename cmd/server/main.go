// Copyright © 2022 Marcus Franke <marcus.franke@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/russross/blackfriday/v2"
	"github.com/spf13/viper"
)

var (
	cmdOut []byte
	mutex  sync.Mutex
)

type yummyConfiguration struct {
	port             string
	repoPath         string
	createrepoBinary string
	rpmBinary        string
	helpFile         string
}

// Main function for the repository webserver
func main() {

	// parse config file
	config := initConfig()
	validateConfig(config)

	repoPath := viper.GetString("yum.repoPath")
	router := httprouter.New()
	router.NotFound = http.FileServer(http.Dir(repoPath))
	router.HandlerFunc("GET", "/help", helpHandler)
	router.HandlerFunc("POST", "/api/upload", apiPostUploadHandler)
	router.HandlerFunc("DELETE", "/api/delete/:filename", apiDeleteHandler)
	log.Fatal(http.ListenAndServe(":"+config.port, router))
}

func updateRepo() bool {
	mutex.Lock()
	workers := viper.GetString("yum.workers")
	repoPath := viper.GetString("yum.repoPath")
	createrepoBinary := viper.GetString("yum.createrepoBinary")
	cmdOut, err := exec.Command(createrepoBinary, "--update", "--workers", workers, repoPath).CombinedOutput()
	if err != nil {
		log.Println(err, string(cmdOut))
		mutex.Unlock()
		return false
	}
	log.Println(string(cmdOut))
	mutex.Unlock()
	return true
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	// get helpFile path from configuration
	helpFile := viper.GetString("yum.helpFile")

	// ingest the configured helpFile
	help, err := ioutil.ReadFile(helpFile)
	if err != nil {
		log.Println("Help file could not be read!")
		http.Error(w, "Could not load the help file", http.StatusInternalServerError)
		return
	}

	// render the Markdown file to HTML using the
	// blackfriday library
	output := blackfriday.Run(help)
	log.Println("/help requested!")
	fmt.Fprintf(w, string(output))
}

func checkAuthentication(r *http.Request) bool {
	// check if authentication is enabled
	if !viper.GetBool("yum.auth.enabled") {
		return true
	}
	// check if header is set
	authString := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authString) != 2 {
		return false
	}
	// decode auth data
	authBytes, err := base64.StdEncoding.DecodeString(authString[1])
	if err != nil {
		return false
	}
	// check user and password
	if string(authBytes) == viper.GetString("yum.auth.user")+":"+viper.GetString("yum.auth.password") {
		return true
	} else {
		return false
	}
}

func apiDeleteHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	filename := params.ByName("filename")
	repoPath := viper.GetString("yum.repoPath")

	if !checkAuthentication(r) {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	if _, err := os.Stat(repoPath + "/" + filename); err == nil {
		// requested file exists
		if err := os.Remove(repoPath + "/" + filename); err != nil {
			errText := fmt.Sprintf("%s - Could not delete file!\n", r.URL)
			log.Printf(errText)
			http.Error(w, errText, http.StatusInternalServerError)
			return
		} else {
			// file deleted
			logText := fmt.Sprintf("%s - File deleted!\n", r.URL.Path)
			log.Printf(logText)
			// update repository
			if !updateRepo() {
				http.Error(w, "Could not update repository", http.StatusInternalServerError)
			}
		}
	} else {
		// file does not exists
		errText := fmt.Sprintf("%s - File not found!\n", r.URL.Path)
		log.Println(errText)
		http.Error(w, errText, http.StatusNotFound)
	}
}

func apiPostUploadHandler(w http.ResponseWriter, r *http.Request) {

	errText := ""
	repoPath := viper.GetString("yum.repoPath")

	file, handler, err := r.FormFile("fileupload")
	if err != nil {
		errText = fmt.Sprintf("%s - incorrect FormFile used, must be fileupload!\n", r.URL)
		log.Println(errText)
		http.Error(w, errText, http.StatusBadRequest)
		return
	}
	defer file.Close()

	if filepath.Ext(handler.Filename) != ".rpm" {
		errText = fmt.Sprintf("%s - %s uploaded, not an rpm package!\n", r.URL, handler.Filename)
		log.Printf(errText)
		http.Error(w, errText, http.StatusUnsupportedMediaType)
		return
	}

	// check if the uploaded file already exists
	// if the repository is configured in protected mode
	// the request will return status 403 (forbidden)
	if viper.GetBool("yum.protected") {
		if _, err := os.Stat(repoPath + "/" + handler.Filename); err == nil {
			errText = fmt.Sprintf("%s - File already exists, forbidden to overwrite!", r.URL)
			log.Println(errText)
			http.Error(w, errText, http.StatusForbidden)
			return
		}
	}

	// create file handler to write uploaded file to
	f, err := os.OpenFile(repoPath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		errText = fmt.Sprintf("%s - %s/%s could not be created!\n", r.URL, repoPath, handler.Filename)
		log.Println(errText)
		http.Error(w, errText, http.StatusInternalServerError)
		log.Fatal(err)
	}
	defer f.Close()

	// copy the file buffer into the file handle
	_, err = io.Copy(f, file)
	if err != nil {
		errText = fmt.Sprintf("%s - an error occurred copying the uploaded file to servers filesystem!\n",
			r.URL)
		log.Println(errText)
		log.Println(err)
		http.Error(w, errText, http.StatusInternalServerError)
		return
	}
	// update repository
	if !updateRepo() {
		http.Error(w, "Could not update repository", http.StatusInternalServerError)
	}
}
