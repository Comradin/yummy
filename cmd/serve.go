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
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	port  string
	debug bool
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the yummy webserver",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("/tmp"))))
		http.HandleFunc("/api/upload", apiuploadhandler)

		log.Fatal(http.ListenAndServe(":8080", nil))
	},
}

func apiuploadhandler(w http.ResponseWriter, r *http.Request) {
	// will handle file uploads
	if debug {
		fmt.Println("Method:", r.Method)
		fmt.Println("Header:", r.Header)
	}

	if r.Method == "POST" {
		file, handler, err := r.FormFile("fileupload")
		if err != nil {
			http.Error(w, "FormFile does not match - use fileupload\n", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// create file handler to write uploaded file to
		f, err := os.OpenFile("/opt/repos/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		io.Copy(f, file)
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
