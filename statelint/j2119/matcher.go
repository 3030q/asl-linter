package j2119

import (
	"fmt"
	"regexp"
	"strings"
)

const must = `(?P<modal>MUST|MAY|MUST NOT)`

var (
	relations = strings.Join([]string{
		"", "equal to", "greater than", "less than",
		"greater than or equal to", "less than or equal to",
	}, "|")
	relation = fmt.Sprintf("((?P<relation>%s)\\s+)", relations)
	// StringMatcher
	s = `"[^"]*"`
	// Non-string value: number, true, false, null
	v          = `\S+`
	relational = fmt.Sprintf("%s(?P<target>%s|%s)", relation, s, v)
	childRole  = `;\s+((its\s+(?P<child_type>value))|` +
		`each\s+(?P<child_type>field|element))` +
		`\s+is\s+an?\s+` +
		`"(?P<child_role>[^"]+)"`
)

type Matcher struct {
	initialized bool
	strings     string
	predicate   string
	roleMatcher string
	roles       []string
	typeRegex   string

	roleDefMatch    *regexp.Regexp
	constraintStart *regexp.Regexp
	constraintMatch *regexp.Regexp
	onlyOneStart    *regexp.Regexp
	onlyOneMatch    *regexp.Regexp
	eachOfMatch     *regexp.Regexp
}

func NewMatcher(root string) *Matcher {
	m := &Matcher{
		roles: make([]string, 0, 1),
	}
	m.AddRole(root)
	m.Constants()
	m.Reconstruct()

	return m
}

func (m *Matcher) Constants() {
	if !m.initialized {
		m.initialized = true
		m.strings = GetOxfordRe(s, OxfordOptions{
			CaptureName:    "strings",
			HasCaptureName: true,
		})
		enum := fmt.Sprintf("one\\s+of\\s+%s", m.strings)
		m.predicate = fmt.Sprintf("(%s|%s)", relational, enum)
	}
}

func (m *Matcher) AddRole(role string) {
	m.roles = append(m.roles, role)
	m.roleMatcher = strings.Join(m.roles, "|")
	m.Reconstruct()
}

func (m *Matcher) Reconstruct() {
	m.MakeTypeRegex()

	// conditional clause
	excludedRoles := "not\\s+" +
		GetOxfordRe(m.roleMatcher, OxfordOptions{
			UseArticle:     true,
			CaptureName:    "excluded",
			HasCaptureName: true,
		}) + "\\s+"

	conditional := "which\\s+is\\s+" + excludedRoles

	// regex for matching constraint lines

	cStart := `^An?\s+` +
		fmt.Sprintf("(?P<role>%s)", m.roleMatcher) + `\s+` +
		fmt.Sprintf("(%s)?", conditional) +
		must + `\s+have\s+an?\s+`

	fieldList := "one\\s+of\\s+" +
		GetOxfordRe(`"[^"]+"`, OxfordOptions{
			CaptureName:    "field_list",
			HasCaptureName: true,
		})

	cMatch := cStart +
		fmt.Sprintf("((?P<type>%s)\\s+)?", m.typeRegex) +
		"field\\s+named\\s+" +
		fmt.Sprintf("((\"(?P<field_name>[^\"]+)\")|(%s))", fieldList) +
		`(\s+whose\s+value\s+MUST\s+be\s+` + m.predicate + ")?" +
		"(" + childRole + ")?" +
		`\.`

	// regexp for matching lines of the form
	// "An X MUST have only one of "Y", "Z", and "W".
	// There's a pattern here, building a separate regex rather than
	// adding more complexity to @constraint_matcher.  Any further
	// additions should be done this way, and
	// TODO: Break @constraint_matcher into a bunch of smaller patterns
	// like this.

	ooStart := `^An?\s+` +
		fmt.Sprintf("(?P<role>%s)", m.roleMatcher) + `\s+` +
		must + `\s+have\s+only\s+`

	ooFieldList := "one\\s+of\\s+" +
		GetOxfordRe(`"[^"]+"`, OxfordOptions{
			CaptureName:    "field_list",
			HasCaptureName: true,
			Connector:      "and",
			HasConnector:   true,
		})

	ooMatch := ooStart + ooFieldList

	// regex for matching role-def lines
	valMatch := "whose\\s+\"(?P<fieldtomatch>[^\"]+)\"" +
		"\\s+field's\\s+value\\s+is\\s+" +
		"(?P<valtomatch>(\"[^\"]*\")|([^\"\\s]\\S+))\\s+"
	withAMatch := "with\\s+an?\\s+\"(?P<with_a_field>[^\"]+)\"\\s+field\\s"

	rdMatch := `^An?\s+` +
		fmt.Sprintf("(?P<role>%s)", m.roleMatcher) + `\s+` +
		fmt.Sprintf("((?P<val_match_present>%s)|(%s))?", valMatch, withAMatch) +
		"is\\s+an?\\s+" +
		"\"(?P<newrole>[^\"]*)\"\\.\\s*$"

	m.roleDefMatch = regexp.MustCompile(rdMatch)

	m.constraintStart = regexp.MustCompile(cStart)
	m.constraintMatch = regexp.MustCompile(cMatch)

	m.onlyOneStart = regexp.MustCompile(ooStart)
	m.onlyOneMatch = regexp.MustCompile(ooMatch)

	eoMatch := "^Each\\s+of\\s" +
		GetOxfordRe(m.roleMatcher, OxfordOptions{
			UseArticle:     true,
			CaptureName:    "each_of",
			HasCaptureName: true,
			Connector:      "and",
			HasConnector:   true,
		}) +
		"\\s+(?P<trailer>.*)$"

	m.eachOfMatch = regexp.MustCompile(eoMatch)
}

// MakeTypeRegex Add modified numeric types to type regex
func (m *Matcher) MakeTypeRegex() {
	types := GetAllTypesString()

	numberTypes := []string{"float", "integer", "numeric"}
	numberModifiers := []string{"positive", "negative", "nonnegative"}

	for _, numberType := range numberTypes {
		for _, modifier := range numberModifiers {
			types = append(types, fmt.Sprintf("%s-%s", modifier, numberType))
		}
	}

	arrayTypes := make([]string, 0, len(types))
	for _, valueType := range types {
		arrayTypes = append(arrayTypes, fmt.Sprintf("%s-array", valueType))
	}

	types = append(types, arrayTypes...)

	nonemptyArrayTypes := make([]string, 0, len(arrayTypes))
	for _, valueType := range arrayTypes {
		nonemptyArrayTypes = append(nonemptyArrayTypes, fmt.Sprintf("nonempty-%s", valueType))
	}

	types = append(types, nonemptyArrayTypes...)

	nullableTypes := make([]string, 0, len(types))
	for _, valueType := range types {
		nullableTypes = append(nullableTypes, fmt.Sprintf("nullable-%s", valueType))
	}

	types = append(types, nullableTypes...)

	m.typeRegex = strings.Join(types, "|")
}

var matchRegex = regexp.MustCompile(`is\s+an?\s+"[^"]*"\.\s*$`)

func (m *Matcher) IsRoleDefLine(line string) bool {
	return matchRegex.MatchString(line)
}

var tokenizeStringsRegex = regexp.MustCompile(`[^"]*"([^"]*)"`)

func (m *Matcher) TokenizeStrings(s string) []string {
	result := make([]string, 0)
	matches := tokenizeStringsRegex.FindAllStringSubmatch(s, -1)

	for _, match := range matches {
		if len(match) > 1 {
			result = append(result, match[1])
		}
	}

	return result
}

var (
	tokenizeValuesMatcher  = regexp.MustCompile(`,|or`)
	splitterTokenizeValues = regexp.MustCompile(`\s+`)
)

func (m *Matcher) TokenizeValues(s string) []string {
	replaced := tokenizeValuesMatcher.ReplaceAllString(s, " ")

	return splitterTokenizeValues.Split(replaced, -1)
}

func (m *Matcher) build(re *regexp.Regexp, line string) map[string]string {
	data := make(map[string]string)
	groupName := re.SubexpNames()
	matches := re.FindAllStringSubmatch(line, -1)

	for _, match := range matches {
		for matchIndex, matchSubstr := range match {
			if groupName[matchIndex] == "" || matchSubstr == "" {
				continue
			}

			data[groupName[matchIndex]] = matchSubstr
		}
	}

	return data
}

func (m *Matcher) BuildRoleDef(line string) map[string]string {
	return m.build(m.roleDefMatch, line)
}

func (m *Matcher) IsConstraintLine(line string) bool {
	return m.constraintStart.MatchString(line)
}

func (m *Matcher) IsOnlyOneMatchLine(line string) bool {
	return m.onlyOneStart.MatchString(line)
}

func (m *Matcher) BuildConstraint(line string) map[string]string {
	return m.build(m.constraintMatch, line)
}

func (m *Matcher) BuildOnlyOne(line string) map[string]string {
	return m.build(m.onlyOneMatch, line)
}

func (m *Matcher) BuildEachOfLine(line string) map[string]string {
	return m.build(m.eachOfMatch, line)
}
