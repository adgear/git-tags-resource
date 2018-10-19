package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adgear/git-tags-resource/actions"
	"github.com/adgear/git-tags-resource/services"
	"github.com/adgear/git-tags-resource/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version = "local"

var (
	rootCmd = &cobra.Command{}

	// Flags
	action      string
	destination string
	source      string
	tmpdir      string
	verbose     bool
	showVersion bool
	help        bool
	stdin       string
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)

	rootCmd.PersistentFlags().StringVarP(&action, "action", "a", "check", "Concourse resource action, defaults to check. Can be check, in or out.")
	rootCmd.PersistentFlags().StringVarP(&destination, "destination", "d", "./", "Destination of input resource")
	rootCmd.PersistentFlags().StringVarP(&source, "source", "s", "./", "Source of output resource")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode.")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false, "Prints the version and exit.")
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "Show help.")
	rootCmd.PersistentFlags().StringVarP(&stdin, "stdin", "S", "", "Feed file instead of stdin")

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	if help {
		rootCmd.Usage()
		os.Exit(0)
	}

	if showVersion {
		log.Info("Version: " + version)
		os.Exit(0)
	}

	var input []byte
	var err error

	if stdin != "" {
		input, err = ioutil.ReadFile(stdin)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else {
		input, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	var inputMap utils.Input

	err = json.Unmarshal(input, &inputMap)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if inputMap.Source.TagFilter == "" {
		log.Fatal("missing tag filter...")
		os.Exit(1)
	}

	gts, err := services.NewGitTagsService(inputMap.Source.PrivateKey, inputMap.Source.PrivateKeyPassword)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	switch action {
	case "check":
		checkResource, err := actions.NewCheckResource(gts)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		output, err := checkResource.Execute(inputMap.Source)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		fmt.Println(output)
		os.Exit(0)
	case "in":
		inResource, err := actions.NewInResource(gts)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		output, err := inResource.Execute(inputMap.Source, destination, inputMap.Version["ref"])
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		fmt.Println(output)
		os.Exit(0)
	default:
		fmt.Println("Nope")
		os.Exit(1)
	}
}
