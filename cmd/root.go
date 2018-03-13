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
	"log"
	"os"
	"os/exec"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "yummy",
	Short: "a simple yum repository server",
	Long: `yummy provides a simple yum repository webserver where you can
upload your rpm packages and download from.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// read configuration
	cobra.OnInitialize(initConfig)

	// Flag for config file
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.yummy.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory
		viper.AddConfigPath(home)

		// Search config parallel to executable
		ex, err := os.Executable()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		exPath := filepath.Dir(ex)

		// Subtract the name of the executable and the last slash
		viper.AddConfigPath(exPath)

		// Search for config file with name ".yummy" (without extension).
		viper.SetConfigName(".yummy")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// validate configuration

	// store configurations in a variables, as we use them quite often
	repoPath := viper.GetString("yum.repopath")
	createrepoBinary := viper.GetString("yum.createrepoBinary")
	rpmBinary := viper.GetString("yum.rpmBinary")
	helpFile := viper.GetString("yum.helpFile")

	// check if repo path exists
	info, err := os.Stat(repoPath)
	if err != nil {
		fmt.Printf("configured repo path '%s' does not exists\n", repoPath)
		os.Exit(1)
	}

	// check if repo path is a directory
	if !info.IsDir() {
		fmt.Printf("configured repo path '%s' is not a directory\n", repoPath)
		os.Exit(1)
	}

	// check if repo path is writeable
	if unix.Access(repoPath, unix.W_OK) != nil {
		fmt.Printf("configured repo path '%s' is not writeable\n", repoPath)
		os.Exit(1)
	}

	// check if createrepo binary exists
	info, err = os.Stat(createrepoBinary)
	if err != nil {
		fmt.Printf("configured createrepo binary '%s' does not exists\n", createrepoBinary)
		os.Exit(1)
	}

	// check if createrepo binary is executable
	if unix.Access(createrepoBinary, unix.X_OK) != nil {
		fmt.Printf("configured createrepo binary '%s' is not executable\n", createrepoBinary)
		os.Exit(1)
	}

	// check if rpm binary exists
	info, err = os.Stat(rpmBinary)
	if err != nil {
		fmt.Printf("configured rpm binary '%s' does not exists\n", rpmBinary)
		os.Exit(1)
	}

	// check if rpm binary is executable
	if unix.Access(rpmBinary, unix.X_OK) != nil {
		fmt.Printf("configured rpm binary '%s' is not executable\n", rpmBinary)
		os.Exit(1)
	}

	// check if help file exists
	info, err = os.Stat(helpFile)
	if err != nil {
		fmt.Printf("configured help file '%s' does not exists\n", helpFile)
		os.Exit(1)
	}

	// initialise repository if not exists
	_, err = os.Stat(repoPath + "/repodata")
	if err != nil {
		fmt.Printf("initialise empty repository in %s\n", repoPath)
		var cmdOut []byte
		cmdOut, err = exec.Command(createrepoBinary, repoPath).CombinedOutput()
		if err != nil {
			log.Println(err, string(cmdOut))
			os.Exit(1)
		}
	} else {
		fmt.Printf("Using existing repository: %s\n", repoPath)
	}
}
