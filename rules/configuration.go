package rules

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/reugn/fsweeper/ospkg"
	"gopkg.in/yaml.v2"
)

const (
	operatorAnd = "AND"
	operatorOr  = "OR"
)

const (
	filterName       = "name"
	filterExtension  = "ext"
	filterSize       = "size"
	filterLastEdited = "lastEdited"
	filterContains   = "contains"
)

const (
	actionEcho   = "echo"
	actionTouch  = "touch"
	actionMove   = "move"
	actionRename = "rename"
	actionDelete = "delete"
)

var (
	// Filters is the supported filters list
	Filters = [...]string{filterName, filterExtension, filterSize, filterLastEdited, filterContains}

	// Actions is the supported actions list
	Actions = [...]string{actionEcho, actionTouch, actionMove, actionRename, actionDelete}
)

// Rule configuration to apply on a file
type Rule struct {
	Path      string   `yaml:"path"`
	Recursive bool     `yaml:"recursive"`
	Operator  string   `yaml:"op" default:"AND"`
	Actions   []Action `yaml:"actions"`
	Filters   []Filter `yaml:"filters"`
}

func (r *Rule) checkFilters(filePath string) bool {
	res := make([]bool, len(r.Filters))
	var i int
	for _, filter := range r.Filters {
		if r.Operator == operatorOr {
			for _, b := range res {
				if b {
					return true
				}
			}
		}

		switch filter.Filter {
		case filterName:
			res[i] = filter.nameFilter(filePath)
		case filterExtension:
			res[i] = filter.extensionFilter(filePath)
		case filterSize:
			res[i] = filter.sizeFilter(filePath)
		case filterLastEdited:
			res[i] = filter.lastEditedFilter(filePath)
		case filterContains:
			res[i] = filter.containsFilter(filePath)
		default:
			log.Fatalf("Unknown filter %s", filter.Filter)
		}

		i++
	}

	return r.filtersResult(res)
}

func (r *Rule) filtersResult(res []bool) bool {
	if r.Operator == operatorOr {
		for _, b := range res {
			if b {
				return true
			}
		}

		return false
	}

	for _, b := range res {
		if !b {
			return false
		}
	}

	return true
}

func (r *Rule) runActions(filePath string, vars *Vars) {
	for _, action := range r.Actions {
		switch action.Action {
		case actionEcho:
			action.echoAction(filePath, vars)
		case actionTouch:
			action.touchAction(filePath)
		case actionMove:
			action.moveAction(filePath, vars)
		case actionRename:
			action.renameAction(filePath, vars)
		case actionDelete:
			action.deleteAction(filePath)
		default:
			log.Fatalf("Unknown action %s", action.Action)
		}
	}
}

// Config is a multiple rules container
type Config struct {
	Vars  Vars   `yaml:"vars"`
	Rules []Rule `yaml:"rules"`
	wg    sync.WaitGroup
}

// ReadConfig reads configuration from the default configuration file
func ReadConfig() *Config {
	return ReadConfigFromFile(GetDefaultConfigFile())
}

// ReadConfigFromFile reads configuration from a custom configuration file
func ReadConfigFromFile(file string) *Config {
	c := &Config{}

	confFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read yaml configuration from file %s, #%v ", file, err)
		return nil
	}

	err = yaml.Unmarshal(confFile, c)
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
		return nil
	}

	// m, _ := yaml.Marshal(c)
	// log.Printf("Parsed configuration file\n%s\n", string(m))

	c.Vars.init()
	return c
}

// ReadConfigFromByteArray reads configuration from a given byte array
func ReadConfigFromByteArray(configYaml []byte) *Config {
	c := &Config{}

	err := yaml.Unmarshal(configYaml, c)
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
		return nil
	}

	c.Vars.init()
	return c
}

// Execute rules
func (c *Config) Execute() {
	for _, rule := range c.Rules {
		c.wg.Add(1)
		go c.iterate(rule)
	}

	c.wg.Wait()
}

func (c *Config) iterate(rule Rule) {
	defer c.wg.Done()

	files, err := ioutil.ReadDir(rule.Path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		filePath := validatePath(rule.Path) + f.Name()
		if f.IsDir() {
			if rule.Recursive {
				rule.Path = filePath
				c.wg.Add(1)
				go c.iterate(rule)
			}
		} else {
			c.wg.Add(1)
			go func(path string) {
				defer c.wg.Done()

				if rule.checkFilters(path) {
					rule.runActions(path, &c.Vars)
				}
			}(filePath)
		}
	}
}

func validatePath(dirPath string) string {
	if strings.HasSuffix(dirPath, ospkg.PathSeparator) {
		return dirPath
	}
	return dirPath + ospkg.PathSeparator
}
