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
	"fmt"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// initConfig parses the config file and ENV variables, if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	// Search config in home directory
	viper.AddConfigPath(home)

	// Search config parallel to executable
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	executablePath := filepath.Dir(executable)
	viper.AddConfigPath(executablePath)

	// Search for config file with name ".yummy" (without extension).
	viper.SetConfigName(".yummy")

	viper.AutomaticEnv() // read in environment variables that match

	// If config file is found, parse it.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// validate configuration

	// store configurations in a variables, as we use them quite often
	repoPath := viper.GetString("yum.repopath")
	createrepoBinary := viper.GetString("yum.createrepoBinary")
	rpmBinary := viper.GetString("yum.rpmBinary")
	helpFile := viper.GetString("yum.helpFile")

	//// check if repo path exists
	//info, err := os.Stat(repoPath)
	//if err != nil {
	//	fmt.Printf("configured repo path '%s' does not exists\n", repoPath)
	//	os.Exit(1)
	//}
	//
	//// check if repo path is a directory
	//if !info.IsDir() {
	//	fmt.Printf("configured repo path '%s' is not a directory\n", repoPath)
	//	os.Exit(1)
	//}
	//
	//// check if repo path is writeable
	//if unix.Access(repoPath, unix.W_OK) != nil {
	//	fmt.Printf("configured repo path '%s' is not writeable\n", repoPath)
	//	os.Exit(1)
	//}

	// check if createrepo binary exists
	//info, err = os.Stat(createrepoBinary)
	//if err != nil {
	//	fmt.Printf("configured createrepo binary '%s' does not exists\n", createrepoBinary)
	//	os.Exit(1)
	//}
	//
	//// check if createrepo binary is executable
	//if unix.Access(createrepoBinary, unix.X_OK) != nil {
	//	fmt.Printf("configured createrepo binary '%s' is not executable\n", createrepoBinary)
	//	os.Exit(1)
	//}
	//
	//// check if rpm binary exists
	//info, err = os.Stat(rpmBinary)
	//if err != nil {
	//	fmt.Printf("configured rpm binary '%s' does not exists\n", rpmBinary)
	//	os.Exit(1)
	//}

	// check if rpm binary is executable
	//if unix.Access(rpmBinary, unix.X_OK) != nil {
	//	fmt.Printf("configured rpm binary '%s' is not executable\n", rpmBinary)
	//	os.Exit(1)
	//}
	//
	//// check if help file exists
	//info, err = os.Stat(helpFile)
	//if err != nil {
	//	fmt.Printf("configured help file '%s' does not exists\n", helpFile)
	//	os.Exit(1)
	//}
	//
	//// initialise repository if not exists
	//_, err = os.Stat(repoPath + "/repodata")
	//if err != nil {
	//	fmt.Printf("initialise empty repository in %s\n", repoPath)
	//	var cmdOut []byte
	//	cmdOut, err = exec.Command(createrepoBinary, repoPath).CombinedOutput()
	//	if err != nil {
	//		log.Println(err, string(cmdOut))
	//		os.Exit(1)
	//	}
	//} else {
	//	fmt.Printf("Using existing repository: %s\n", repoPath)
	//}
}
