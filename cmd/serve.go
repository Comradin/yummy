// Copyright © 2017 Marcus Franke <marcus.franke@gmail.com>
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

package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/russross/blackfriday.v2"
)

var (
	port   string
	cmdOut []byte
	mutex  sync.Mutex
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the yummy webserver",
	Run: func(cmd *cobra.Command, args []string) {

		repoPath := viper.GetString("yum.repopath")

		router := httprouter.New()
		router.Handler("GET", "/", http.FileServer(http.Dir(repoPath)))

		router.GET("/help", helpHandler)
		router.POST("/api/upload", apiPostUploadHandler)
		//router.PUT("/api/upload/:filename", apiUploadPut)
		//router.DELETE("/api/delete/:name", apiDeleteHandler)

		log.Fatal(http.ListenAndServe(":8080", router))
	},
}

func helpHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func apiPostUploadHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	errText := ""
	repoPath := viper.GetString("yum.repopath")
	workers := viper.GetString("yum.workers")
	createrepoBinary := viper.GetString("yum.createrepoBinary")

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
			errText = fmt.Sprintf("%s - File already exists, forbidden to overwrite!\n", r.URL)
			log.Println(errText)
			log.Println(err)
			http.Error(w, errText, http.StatusForbidden)
			return
		}
		log.Println("File already exists, will overwrite: " + handler.Filename)
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
		errText = fmt.Sprintf("%s - an error occured copying the uploaded file to servers filesystem!\n",
			r.URL)
		log.Println(errText)
		log.Println(err)
		http.Error(w, errText, http.StatusInternalServerError)
		return
	}

	// process the uploaded file
	mutex.Lock()
	cmdOut, err = exec.Command(createrepoBinary, "--update", "--workers", workers, repoPath).CombinedOutput()
	if err != nil {
		fmt.Fprintln(w, string(cmdOut))
		http.Error(w, "Could not update repository", http.StatusInternalServerError)
		log.Println(err, string(cmdOut))
		mutex.Unlock()
		return
	}
	log.Println(string(cmdOut))
	mutex.Unlock()
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Flags for the serve command.
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
}
