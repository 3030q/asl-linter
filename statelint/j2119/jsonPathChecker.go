package j2119

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	initialNameClasses   = []string{"Lu", "Ll", "Lt", "Lm", "Lo", "Nl"}
	followingNameClasses = []string{"Lu", "Ll", "Lt", "Lm", "Lo", "Nl", "Mn", "Mc", "Nd", "Pc"}
	dotSeparator         = `\.\.?`

	nameRe        = classesToRe(initialNameClasses) + classesToRe(followingNameClasses) + "*"
	dotStep       = fmt.Sprintf("%s((%s)|(\\*))", dotSeparator, nameRe)
	rpDotStep     = dotSeparator + nameRe
	bracketStep   = fmt.Sprintf("\\['%s'\\]", nameRe)
	rpNumIndex    = `\[\d+\]`
	numIndex      = `\[\d+(, *\d+)?\]`
	starIndex     = `\[\*\]`
	colonIndex    = `\[(-?\d+)?:(-?\d+)?\]`
	index         = fmt.Sprintf("((%s)|(%s)|(%s))?", numIndex, starIndex, colonIndex)
	step          = fmt.Sprintf("((%s)|(%s)|(%s))(%s)?", dotStep, bracketStep, index, index)
	rpStep        = fmt.Sprintf("((%s)|(%s))(%s)?", rpDotStep, bracketStep, rpNumIndex)
	path          = fmt.Sprintf("^\\$(%s)*$", step)
	referencePath = fmt.Sprintf("^\\$(%s)*$", rpStep)

	pathRe          = regexp.MustCompile(path)
	referencePathRe = regexp.MustCompile(referencePath)
)

func classesToRe(classes []string) string {
	newClasses := make([]string, 0, len(classes))
	for _, class := range classes {
		newClasses = append(newClasses, fmt.Sprintf("\\p{%s}", class))
	}

	return fmt.Sprintf("[%s]", strings.Join(newClasses, ""))
}

func IsPath(path string) bool {
	return pathRe.MatchString(path)
}

func IsReferencePath(path string) bool {
	return referencePathRe.MatchString(path)
}
