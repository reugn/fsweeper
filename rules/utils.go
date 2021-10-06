package rules

import "os"

// DefaultConfigFile is the configuration file default path.
const DefaultConfigFile string = "conf.yaml"

// GetDefaultConfigFile returns a configuration file path.
func GetDefaultConfigFile() string {
	confFile, ok := os.LookupEnv("FSWEEPER_CONFIG_FILE")
	if !ok {
		confFile = DefaultConfigFile
	}
	return confFile
}
