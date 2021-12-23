package j2119

import (
	"fmt"
	"regexp"
	"strconv"
)

const float64Bits = 64

var (
	stringMatcher = regexp.MustCompile(`^"(.*)"$`)
	intMatcher    = regexp.MustCompile(`^-?\d+$`)
	floatMatcher  = regexp.MustCompile(`^-?\d+.?\d+?$`)
)

func DeduceValue(value string) interface{} {
	switch {
	case stringMatcher.MatchString(value):
		return stringMatcher.FindStringSubmatch(value)[1]
	case value == "true":
		return true
	case value == "false":
		return false
	case value == "null":
		return nil
	case intMatcher.MatchString(value):
		integer, err := strconv.Atoi(value)
		if err != nil {
			panic("err when atoi")
		}

		return integer
	case floatMatcher.MatchString(value):
		float, err := strconv.ParseFloat(value, float64Bits)
		if err != nil {
			panic(fmt.Sprintf("err when parse float %s", value))
		}

		return float
	default:
		return value
	}
}
