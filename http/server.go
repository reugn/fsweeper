package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/reugn/fsweeper/rules"
)

// StartHTTPServer starts HTTP listener
func StartHTTPServer(host string, port int) {
	addr := host + ":" + strconv.Itoa(port)

	http.HandleFunc("/", rootActionHandler)

	// executes configuration
	http.HandleFunc("/execute", executeActionHandler)

	// writes rules to the default configuration file
	http.HandleFunc("/config", configActionHandler)

	http.ListenAndServe(addr, nil)
}

func rootActionHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprintf(w, "")
}

func executeActionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		config := rules.ReadConfig()
		config.Execute()
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		config := rules.ReadConfigFromByteArray(body)
		config.Execute()
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func configActionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	err = ioutil.WriteFile(rules.GetDefaultConfigFile(), body, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
