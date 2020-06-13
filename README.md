# fsweeper
*The file management automation tool.*  

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
`fsweeper` uses YAML as a configuration file format.  
A default configuration file name could be set using the `FSWEEPER_CONFIG_FILE` environment variable. It will fallback to `conf.yaml` otherwise.  
Use `--conf=custom.yaml` parameter to run against a custom configuration file.  

Let's start with a simple example:
```yaml
rules:
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
This configuration will look for files that are bigger than 10000 bytes or empty and print out "Found size" message on each.

What if we want to move all JSON files to the archive folder:
```yaml
rules:
    - path: ./examples/files
      recursive: true
      actions:
        - action: move
          payload: ./examples/archive
        - action: echo
          payload: "Found json file"
      filters:
        - filter: ext
          payload: .json
```

### Variables and pipelines
We can build dynamic configurations using variables and pipelines.
Variables introduce a set of dynamic placeholders to be substituted in runtime.
Pipelines chain together a series of template commands to compactly express a series of transformations.

Let's see an example:
```yaml
vars:
    dateTimeFormat: "2006.01.02[15-04-05]"
rules:
    - path: ./examples/files
      recursive: true
      op: AND
      actions:
        - action: echo
          payload: "Found {{ .FileName | upper | quote }} at {{ .DateTime }}, size {{ .FileSize }}"
      filters:
        - filter: name
          payload: "[0-9]+"
        - filter: contains
          payload: "1234"
```

| Variables | Description |
| --- | --- |
| .FileSize | Returns context file size |
| .FileName | Returns context file name |
| .FilePath | Returns context file full path |
| .FileExt | Returns context file extension |
| .Time | Time default format: "15-04-05" | 
| .Date | Date default format: "2006-01-02" |
| .DateTime | DateTime default format: "2006-01-02[15-04-05]" |
| .Ts | Unix epoch time |

You can override `timeFormat`, `dateFormat`, and `dateTimeFormat` in the configuration file.

| Pipeline functions | Description |
| --- | --- |
| trim | Removes leading and trailing white space |
| upper | Returns string with all Unicode letters mapped to their upper case |
| lower | Returns string with all Unicode letters mapped to their lower case |
| quote | Wraps string with double quotes |
| head n | Takes first "n" characters of the input string |
| len | Returns string representation of an input length |

## License
MIT