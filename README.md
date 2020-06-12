# fsweeper
The file management automation tool.  
Got tired of cleaning and organizing the disk time after time? Get a couple of minutes to read through this guide and write rules to automate the process.

## Installation
Go installation required to build `fsweeper`.  
To install the latest stable version of Go, visit [http://golang.org/dl/](http://golang.org/dl/).
```
go build
```

## Usage
`fsweeper` can be used in either a server or CLI mode.  
To run it in a server mode: `./fsweeper --http`.  
Running in a server mode will expose the following HTTP endpoints:
* `GET /` - health endpoint
* `GET /execute` - execute rules from the default configuration file
* `POST /execute` - execute a provided yaml configuration
* `POST /config` - write rules to the default configuration file

Command line options list:
```
./fsweeper --help

Usage of ./fsweeper:
  -actions
        Show supported actions
  -conf string
        Configuration file path (default "conf.yaml")
  -configure
        Open default configuration file in $EDITOR
  -filters
        Show supported filters
  -host string
        HTTP Server host (default "0.0.0.0")
  -http
        Run in a HTTP mode
  -port int
        HTTP Server port (default 8081)
  -version
        Show version
```

## Configuration
A default configuration file name could be set using the `FSWEEPER_CONFIG_FILE` environment variable. It will fallback to `conf.yaml` otherwise.

Use `--conf=custom.yaml` parameter to run against a custom configuration file.

Configuration file example:
```yaml
rules:
    - path: ./examples/files
      recursive: true
      op: AND
      actions:
        - action: touch
        - action: echo
          payload: "Found name && contains"
      filters:
        - filter: name
          payload: "[0-9]+"
        - filter: contains
          payload: "1234"

    - path: ./examples/files
      recursive: true
      actions:
        - action: rename
          payload: foo.json
        - action: echo
          payload: "Found ext"
      filters:
        - filter: ext
          payload: .json

    - path: ./examples/files
      recursive: true
      op: OR
      actions:
        - action: echo
          payload: "Found size"
      filters:
        - filter: size
          payload: gt 10000
        - filter: size
          payload: eq 0
```

## License
MIT