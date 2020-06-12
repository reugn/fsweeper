package rules

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Filter to verify on a file
type Filter struct {
	Filter  string `yaml:"filter"`
	Payload string `yaml:"payload"`
}

func (f *Filter) nameFilter(filePath string) bool {
	regex := regexp.MustCompile(f.Payload)
	return regex.MatchString(filePath)
}

func (f *Filter) extensionFilter(filePath string) bool {
	return strings.HasSuffix(filePath, f.Payload)
}

func (f *Filter) containsFilter(filePath string) bool {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read file %s. %s\n", filePath, err.Error)
		return false
	}

	regex := regexp.MustCompile(f.Payload)
	return regex.MatchString(string(data))
}

func (f *Filter) sizeFilter(filePath string) bool {
	p := strings.Split(f.Payload, " ")
	if len(p) != 2 {
		log.Fatalf("Invalid payload for sizeFilter %s", f.Payload)
	}

	size, err := strconv.ParseInt(p[1], 10, 64)
	if err != nil {
		log.Fatalf("Invalid filter size configuration %s", p[1])
	}

	fi, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("Failed to get %s Stat()", filePath)
	}

	if lt(p[0]) {
		return fi.Size() < size
	} else if eq(p[0]) {
		return fi.Size() == size
	}

	// default gt operator
	return fi.Size() > size
}

func (f *Filter) lastEditedFilter(filePath string) bool {
	p := strings.Split(f.Payload, " ")
	if len(p) != 2 {
		log.Fatalf("Invalid payload for lastEditedFilter %s", f.Payload)
	}

	ts, err := strconv.ParseInt(p[1], 10, 64)
	if err != nil {
		log.Fatalf("Invalid last edited configuration %s", p[1])
	}

	fi, err := os.Stat(filePath)
	if err != nil {
		log.Fatalf("Failed to get %s Stat()", filePath)
	}

	if lt(p[0]) {
		return fi.ModTime().Before(time.Unix(ts, 0))
	} else if eq(p[0]) {
		return fi.ModTime().Equal(time.Unix(ts, 0))
	}

	// default gt operator
	return fi.ModTime().After(time.Unix(ts, 0))
}

func gt(op string) bool {
	return op == "gt" || op == ">"
}

func lt(op string) bool {
	return op == "lt" || op == "<"
}

func eq(op string) bool {
	return op == "eq" || op == "=" || op == "=="
}
