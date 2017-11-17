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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/russross/blackfriday.v2"
	"sync"
)

var (
	port   string
	debug  bool
	cmdOut []byte
	mutex  sync.Mutex
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the yummy webserver",
	Run: func(cmd *cobra.Command, args []string) {

		repoPath := viper.GetString("yum.repopath")
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(repoPath))))
		http.HandleFunc("/help", helpHandler)
		http.HandleFunc("/api/upload", apiuploadhandler)

		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	// get helpFile path from configuration
	helpFile := viper.GetString("yum.helpFile")

	// ingest the configured helpFile
	help, err := ioutil.ReadFile(helpFile)
	if err != nil {
		http.Error(w, "Could not load the help file", http.StatusInternalServerError)
		return
	}

	// render the Markdown file to HTML using the
	// blackfriday library
	output := blackfriday.Run(help)
	fmt.Fprintf(w, string(output))
}

func apiuploadhandler(w http.ResponseWriter, r *http.Request) {

	repoPath := viper.GetString("yum.repopath")
	workers := viper.GetString("yum.workers")
	createrepoBinary := viper.GetString("yum.createrepoBinary")

	if debug {
		fmt.Println("Method:", r.Method)
		fmt.Println("Header:", r.Header)
		fmt.Println("repoPath:", repoPath)
	}

	// will handle file uploads
	if r.Method == "POST" {
		file, handler, err := r.FormFile("fileupload")
		if err != nil {
			http.Error(w, "FormFile does not match - use fileupload\n", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if filepath.Ext(handler.Filename) != ".rpm" {
			http.Error(w, "File not RPM\n", http.StatusUnsupportedMediaType)
			return
		}

		// create file handler to write uploaded file to
		f, err := os.OpenFile(repoPath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			http.Error(w, "An error occurred", http.StatusInternalServerError)
			log.Fatal(err)
		}
		defer f.Close()

		// copy the file buffer into the file handle
		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, "An error occurred applying the upload to the filesystem", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		// process the uploaded file
		mutex.Lock()
		cmdOut, err = exec.Command(createrepoBinary, "-v", "-p", "--update", "--workers", workers, repoPath).Output()
		if err != nil {
			http.Error(w, "Could not update repository", http.StatusInternalServerError)
			log.Println(err)
			mutex.Unlock()
			return
		}
		mutex.Unlock()

		if debug {
			log.Println(string(cmdOut))
		}
	}

	// assume curl --upload-file style of upload type
	// this is currently not supported
	if r.Method == "PUT" {
		http.Error(w, "Method not allowed, POST binary to URI\n", http.StatusMethodNotAllowed)
	}
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Flags for the serve command.
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to listen on")
	serveCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")
}
