# FSweeper
<img src="docs/images/fsweeper.png" align='right'/>

[![Build](https://github.com/reugn/fsweeper/actions/workflows/build.yml/badge.svg)](https://github.com/reugn/fsweeper/actions/workflows/build.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/reugn/fsweeper)](https://pkg.go.dev/github.com/reugn/fsweeper)
[![Go Report Card](https://goreportcard.com/badge/github.com/reugn/fsweeper)](https://goreportcard.com/report/github.com/reugn/fsweeper)

An intuitive and simple file management automation tool.  
Read this guide and write rules for organizing file storage in just a couple of minutes.

## Installation
Pick a binary from the [releases](https://github.com/reugn/fsweeper/releases).

### Build from source 
Download and install Go https://golang.org/doc/install.

Get the package:
```sh
go get github.com/reugn/fsweeper
```

Read this [guide](https://golang.org/doc/tutorial/compile-install) on how to compile and install the application.

## Usage
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

What if we want to move all JSON files to an archive folder:
```yaml
rules:
    - path: ./examples/files
      recursive: true
      actions:
        - action: move
          payload: ./examples/archive
        - action: echo
          payload: "Found JSON file"
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

| Variables   | Description |
| ----------- | --- |
| `.FileSize` | Returns the context file size |
| `.FileName` | Returns the context file name |
| `.FilePath` | Returns the context file full path |
| `.FileExt`  | Returns the context file extension |
| `.Time`     | Current Time (default format: "15-04-05")<sup>[1](#format)</sup> | 
| `.Date`     | Current Date (default format: "2006-01-02")<sup>[1](#format)</sup> |
| `.DateTime` | Current DateTime (default format: "2006-01-02[15-04-05]")<sup>[1](#format)</sup> |
| `.Ts`       | Current Unix epoch time |

<sup name="format">1</sup>You can override `timeFormat`, `dateFormat`, and `dateTimeFormat` in the configuration file.

| Pipeline functions | Description |
| -------- | --- |
| `trim`   | Removes leading and trailing white spaces |
| `upper`  | Returns a string with all Unicode letters mapped to their upper case |
| `lower`  | Returns a string with all Unicode letters mapped to their lower case |
| `quote`  | Wraps a string with double quotes |
| `head n` | Takes the first "n" characters of the input string |
| `len`    | Returns a string representation of the input length |

## License
MIT
