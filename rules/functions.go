package rules

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	trimFunc  = "trim"
	upperFunc = "upper"
	lowerFunc = "lower"
	quoteFunc = "quote"
	headFunc  = "head"
	lenFunc   = "len"
)

// Pipeline chains together series of template commands to compactly express
// series of transformations.
func Pipeline(in string, chain []string) string {
	out := in

	for _, f := range chain {
		switch strings.ToLower(f) {
		case trimFunc:
			out = strings.TrimSpace(out)
		case upperFunc:
			out = strings.ToUpper(out)
		case lowerFunc:
			out = strings.ToLower(out)
		case quoteFunc:
			out = fmt.Sprintf("\"%s\"", out)
		case lenFunc:
			out = strconv.Itoa(len(out))
		}

		if strings.HasPrefix(f, headFunc) {
			tkn := strings.Split(f, " ")
			if len(tkn) != 2 {
				log.Fatalf("Invalid len function configuration: %s", f)
			}
			out = out[:len(tkn[1])]
		}
	}

	return out
}
