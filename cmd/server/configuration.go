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
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

func init() {
	flag.StringVar(&cfgFile, "config", "$HOME/.yummy.yaml", "config file (default is $HOME/.yummy.yaml)")
	flag.Parse()
}

// initConfig parses the config file and ENV variables, if set.
func initConfig() yummyConfiguration {

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

	if cfgFile != "$HOME/.yummy.yaml" {
		// Use config file from the flag.
		fmt.Println("Checking config file:", cfgFile)
		viper.SetConfigFile(cfgFile)
	}

	// If config file is found, parse it.
	if err = viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// store configurations in a variables, as we use them quite often
	return yummyConfiguration{
		repopath:         viper.GetString("yum.repopath"),
		createrepoBinary: viper.GetString("yum.createrepoBinary"),
		rpmBinary:        viper.GetString("yum.rpmBinary"),
		helpFile:         viper.GetString("yum.helpFile"),
		port:             viper.GetString("yum.port"),
	}
}
