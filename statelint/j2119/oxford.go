package j2119

import (
	"fmt"
	"regexp"
	"strings"
)

// BASIC - (?P<CAPTURE>X((((,\\s+X)+,)?)?\\s+or\\s+X)?)

const (
	basic        = `(?P<CAPTURE>X((((,\s+X)+,)?)?\s+or\s+X)?)`
	hasCapture   = `(?P<CAPTURE>`
	inter        = `((((,\s+`
	hasConnector = `)+,)?)?\s+or\s+`
	last         = `)?)`
)

type OxfordOptions struct {
	UseArticle     bool
	CaptureName    string
	HasCaptureName bool
	Connector      string
	HasConnector   bool
}

func GetOxfordRe(particle string, opts OxfordOptions) string {
	var lHasCapture, lInter, lHasConnector, lLast = hasCapture, inter, hasConnector, last

	if opts.HasConnector {
		replacer := regexp.MustCompile("or")
		lHasConnector = replacer.ReplaceAllString(lHasConnector, opts.Connector)
	}

	if opts.UseArticle {
		particle = fmt.Sprintf("an?\\s+(%s)", particle)
	} else {
		particle = fmt.Sprintf("(%s)", particle)
	}

	if opts.HasCaptureName {
		captureMatcher := regexp.MustCompile("CAPTURE")
		lHasCapture = captureMatcher.ReplaceAllString(lHasCapture, opts.CaptureName)
	} else {
		captureMatcher := regexp.MustCompile(`\?P<CAPTURE>`)
		lHasCapture = captureMatcher.ReplaceAllString(lHasCapture, "")
	}

	return strings.Join([]string{lHasCapture, lInter, lHasConnector, lLast}, particle)
}

var breakStringListRegex = regexp.MustCompile(`[^"]*"([^"]*)"`)

func BreakStringList(list string) []string {
	result := make([]string, 0)

	for _, match := range breakStringListRegex.FindAllStringSubmatch(list, -1) {
		result = append(result, match[1])
	}

	return result
}

func BreakRoleList(matcher Matcher, list string) []string {
	pieces := make([]string, 0)
	re := regexp.MustCompile(fmt.Sprintf(`an?\s+(%s)(,\s+)?`, matcher.roleMatcher))

	for _, match := range re.FindAllStringSubmatch(list, -1) {
		pieces = append(pieces, match[1])
	}

	list = re.ReplaceAllString(list, "")
	re = regexp.MustCompile(fmt.Sprintf(`\s*(and|or)\s+an?\s+(%s)`, matcher.roleMatcher))

	if re.MatchString(list) {
		pieces = append(pieces, re.FindStringSubmatch(list)[2])
	}

	return pieces
}
