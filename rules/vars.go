package rules

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/reugn/fsweeper/ospkg"
	"gopkg.in/yaml.v2"
)

const (
	defaultTimeFormat     = "15-04-05"
	defaultDateFormat     = "2006-01-02"
	defaultDateTimeFormat = "2006-01-02[15-04-05]"
)

const (
	fileSizeVar = ".FileSize"
	fileNameVar = ".FileName"
	filePathVar = ".FilePath"
	fileExtVar  = ".FileExt"
	timeVar     = ".Time"
	dateVar     = ".Date"
	dateTimeVar = ".DateTime"
	tsVar       = ".Ts"
)

var re *regexp.Regexp = regexp.MustCompile(`\{\{(.+?)\}\}`)

// Vars are the configuration variables to substitute
type Vars struct {
	TimeFormat     string `yaml:"timeFormat"`
	DateFormat     string `yaml:"dateFormat"`
	DateTimeFormat string `yaml:"dateTimeFormat"`
	now            time.Time
}

// String is a Stringer interface implementation
func (v *Vars) String() string {
	arr, _ := yaml.Marshal(*v)
	return string(arr)
}

// Process payload string
func (v *Vars) Process(str string, ctx string) string {
	compile := func(s string) string {
		return v.processVariableBlock(s, ctx)
	}
	return re.ReplaceAllStringFunc(str, compile)
}

// compile single variable block
func (v *Vars) processVariableBlock(variable string, ctx string) string {
	variable = strings.Trim(variable, "{{")
	variable = strings.Trim(variable, "}}")
	tokens := strings.Split(variable, "|")

	// trim spaces
	for i := range tokens {
		tokens[i] = strings.TrimSpace(tokens[i])
	}

	variable, chain := tokens[0], tokens[1:]

	var value string
	switch variable {
	case fileSizeVar:
		value = getFileSize(ctx)
	case fileNameVar:
		value = getFileName(ctx)
	case filePathVar:
		value = ctx
	case fileExtVar:
		value = getFileExtension(ctx)
	case timeVar:
		value = v.getTime()
	case dateVar:
		value = v.getDate()
	case dateTimeVar:
		value = v.getDateTime()
	case tsVar:
		value = strconv.FormatInt(v.now.Unix(), 10)
	default:
		log.Fatalf("Unknown variable configuration: %s", variable)
	}

	return Pipeline(value, chain)
}

// init variables, set default value if not set
func (v *Vars) init() {
	if v.TimeFormat == "" {
		v.TimeFormat = defaultTimeFormat
	}

	if v.DateFormat == "" {
		v.DateFormat = defaultDateFormat
	}

	if v.DateTimeFormat == "" {
		v.DateTimeFormat = defaultDateTimeFormat
	}

	v.now = time.Now()
}

// returns file size
func getFileSize(filePath string) string {
	fi, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("Failed to get %s Stat()", filePath)
	}
	return strconv.FormatInt(fi.Size(), 10)
}

// returns file name
func getFileName(filePath string) string {
	p := strings.Split(filePath, ospkg.PathSeparator)
	return p[len(p)-1]
}

// returns file extension
func getFileExtension(filePath string) string {
	p := strings.Split(filePath, ".")
	return p[len(p)-1]
}

// returns formatted time
func (v *Vars) getTime() string {
	return v.now.Format(v.TimeFormat)
}

// returns formatted date
func (v *Vars) getDate() string {
	return v.now.Format(v.DateFormat)
}

// returns formatted datetime
func (v *Vars) getDateTime() string {
	return v.now.Format(v.DateTimeFormat)
}
