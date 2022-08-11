package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func validateConfig(config yummyConfiguration) {

	// check if repo path exists
	info, err := os.Stat(config.repoPath)
	if err != nil {
		log.Fatalf("configured repo path '%s' does not exists\n", config.repoPath)
	}

	// check if repo path is a directory
	if !info.IsDir() {
		log.Fatalf("configured repo path '%s' is not a directory\n", config.repoPath)
	}

	// check if repo path is writeable
	if unix.Access(config.repoPath, unix.W_OK) != nil {
		log.Fatalf("configured repo path '%s' is not writeable\n", config.repoPath)
	}

	//check if createrepo binary exists
	info, err = os.Stat(config.createrepoBinary)
	if err != nil {
		log.Fatalf("configured createrepo binary '%s' does not exists\n", config.createrepoBinary)
	}

	// check if createrepo binary is executable
	if unix.Access(config.createrepoBinary, unix.X_OK) != nil {
		log.Fatalf("configured createrepo binary '%s' is not executable\n", config.createrepoBinary)
	}

	// check if rpm binary exists
	info, err = os.Stat(config.rpmBinary)
	if err != nil {
		log.Fatalf("configured rpm binary '%s' does not exists\n", config.rpmBinary)
	}

	//check if rpm binary is executable
	if unix.Access(config.rpmBinary, unix.X_OK) != nil {
		log.Fatalf("configured rpm binary '%s' is not executable\n", config.rpmBinary)
	}

	// check if help file exists
	info, err = os.Stat(config.helpFile)
	if err != nil {
		log.Fatalf("configured help file '%s' does not exists\n", config.helpFile)
	}

	// initialise repository if not exists
	_, err = os.Stat(config.repoPath + "/repodata")
	if err != nil {
		fmt.Printf("initialise empty repository in %s\n", config.repoPath)
		var cmdOut []byte
		cmdOut, err = exec.Command(config.createrepoBinary, config.repoPath).CombinedOutput()
		if err != nil {
			log.Fatal(err, string(cmdOut))
		}
	} else {
		fmt.Printf("Using existing repository: %s\n", config.repoPath)
	}
}
