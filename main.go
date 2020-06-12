package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/reugn/fsweeper/http"
	"github.com/reugn/fsweeper/rules"
)

var version = "0.1.0"

var (
	configFile      = flag.String("conf", rules.DefaultConfigFile, "Configuration file path")
	httpParam       = flag.Bool("http", false, "Run in a HTTP mode")
	serverHostParam = flag.String("host", "0.0.0.0", "HTTP Server host")
	serverPortParam = flag.Int("port", 8081, "HTTP Server port")

	configureParam = flag.Bool("configure", false, "Open default configuration file in $EDITOR")
	versionParam   = flag.Bool("version", false, "Show version")
	filtersParam   = flag.Bool("filters", false, "Show supported filters")
	actionsParam   = flag.Bool("actions", false, "Show supported actions")
)

func main() {
	flag.Parse()

	if *configureParam {
		openFileInEditor(rules.GetDefaultConfigFile())
		return
	}

	if rt := handleInfoFlags(); rt {
		return
	}

	if *httpParam {
		log.Printf("Starting HTTP server on port %d...\n", *serverPortParam)
		http.StartHTTPServer(*serverHostParam, *serverPortParam)
	} else {
		// read configuration file
		config := rules.ReadConfigFromFile(*configFile)

		// execute rules
		log.Println("Starting execute rules...")
		start := time.Now()
		config.Execute()

		log.Printf("Done in %v.\n", time.Since(start))
	}
}

func handleInfoFlags() bool {
	var rt bool

	if *versionParam {
		fmt.Println("Version: " + version)
		rt = true
	}

	if *filtersParam {
		fmt.Println("Filters:")
		for _, filter := range rules.Filters {
			fmt.Printf("\t- %s\n", filter)
		}
		rt = true
	}

	if *actionsParam {
		fmt.Println("Actions:")
		for _, action := range rules.Actions {
			fmt.Printf("\t- %s\n", action)
		}
		rt = true
	}

	return rt
}

func openFileInEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	// get the full executable path for the editor
	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
