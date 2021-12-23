package j2119

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var roles = []string{
	"Pass State", "Task State", "Choice State",
	"Parallel State", "Succeed State", "Fail State", "Task Tate",
}

var eachOfLines = []string{
	"Each of a Pass State, a Task State, a Choice State, and a Parallel State MAY have a boolean field named \"End\".",
	"Each of a Succeed State and a Fail State is a \"Terminal State\".",
	"Each of a Task State and a Parallel State MAY have an object-array field named \"Catch\"; each member is a \"Catcher\".",
}

func TestMatcher_SpotEachOfLines(t *testing.T) {
	t.Parallel()

	m := NewMatcher("message")
	for _, role := range roles {
		m.AddRole(role)
	}

	for _, line := range eachOfLines {
		assert.True(t, m.eachOfMatch.MatchString(line))
	}
}

func TestMatcher_HandleOnlyOneOfLine(t *testing.T) {
	t.Parallel()

	line := `A x MUST have only one of "Seconds", "SecondsPath", "Timestamp", and "TimestampPath".`
	m := NewMatcher("x")
	assert.True(t, m.IsOnlyOneMatchLine(line))

	match := m.BuildOnlyOne(line)
	role, ok := match["role"]

	assert.True(t, ok)
	assert.Equal(t, role, "x")

	s := match["field_list"]
	l := BreakStringList(s)

	for _, p := range []string{"Seconds", "SecondsPath", "Timestamp", "TimestampPath"} {
		has := false

		for i := range l {
			if l[i] == p {
				has = true

				break
			}
		}

		if !has {
			assert.Fail(t, fmt.Sprintf("can not find %s in %v", p, l))
		}
	}
}

var rdLines = []string{
	"A State whose \"End\" field's value is true is a \"Terminal State\".",
	"Each of a Succeed State and a Fail state is a \"Terminal State\".",
	"A Choice Rule with a \"Variable\" field is a \"Comparison\".",
}

func TestMatcher_SpotRoleDefLines(t *testing.T) {
	t.Parallel()

	m := NewMatcher("message")
	for _, line := range rdLines {
		assert.True(t, m.IsRoleDefLine(line))
	}
}

var valueBasedRoleDefs = []string{
	"A State whose \"End\" field's value is true is a \"Terminal State\".",
	"A State whose \"Comment\" field's value is \"Hi\" is a \"Frobble\".",
	"A State with a \"Foo\" field is a \"Bar\".",
}

func TestMatcher_MatchValueBasedRoleDefs(t *testing.T) {
	t.Parallel()

	m := NewMatcher("State")

	for _, line := range valueBasedRoleDefs {
		assert.True(t, m.roleDefMatch.MatchString(line))
	}

	fields := m.BuildRoleDef(valueBasedRoleDefs[0])
	assert.True(t, fields["role"] == "State")
	assert.True(t, fields["fieldtomatch"] == "End")
	assert.True(t, fields["valtomatch"] == "true")
	assert.True(t, fields["newrole"] == "Terminal State")
	_, has := fields["val_match_present"]
	assert.True(t, has)

	fields = m.BuildRoleDef(valueBasedRoleDefs[1])
	assert.True(t, fields["role"] == "State")
	assert.True(t, fields["fieldtomatch"] == "Comment")
	assert.True(t, fields["valtomatch"] == "\"Hi\"")
	assert.True(t, fields["newrole"] == "Frobble")
	_, has = fields["val_match_present"]
	assert.True(t, has)

	fields = m.BuildRoleDef(valueBasedRoleDefs[2])
	assert.True(t, fields["role"] == "State")
	assert.True(t, fields["newrole"] == "Bar")
	_, has = fields["with_a_field"]
	assert.True(t, has)
}

func TestMatcher_MatchIsARoleDefs(t *testing.T) {
	t.Parallel()

	m := NewMatcher("Foo")
	assert.True(t, m.roleDefMatch.MatchString(`A Foo is a "Bar".`))
}

func TestMatcher_ProperlyParseIsARoleDefs(t *testing.T) {
	t.Parallel()

	m := NewMatcher("Foo")
	m.AddRole("Bar")
	f := m.BuildRoleDef(`A Foo is a "Bar".`)
	_, has := f["val_match_present"]
	assert.True(t, !has)
}

var lines = []string{
	`A message MUST have an object field named "States"; each field is a "State".`,
	`A message MUST have a negative-integer-array field named "StartAt".`,
	`A message MAY have a string-array field named "StartAt".`,
	`A message MUST NOT have a field named "StartAt".`,
	`A message MUST have a field named one of "StringEquals", "StringLessThan", "StringGreaterThan", "StringLessThanEquals", "StringGreaterThanEquals", "NumericEquals", "NumericLessThan", "NumericGreaterThan", "NumericLessThanEquals", "NumericGreaterThanEquals", "BooleanEquals", "TimestampEquals", "TimestampLessThan", "TimestampGreaterThan", "TimestampLessThanEquals", or "TimestampGreaterThanEquals".`,
}

func TestMatcher_SpotSimpleConstraintLine(t *testing.T) {
	t.Parallel()

	m := NewMatcher("message")
	for _, line := range lines {
		assert.True(t, m.IsConstraintLine(line))
	}
}

func TestMatcher_SpotSimpleConstraintLineWithNewRoles(t *testing.T) {
	t.Parallel()

	m := NewMatcher("message")
	m.AddRole("avatar")

	for _, line := range lines {
		newLine := strings.ReplaceAll(line, "message", "avatar")
		assert.True(t, m.IsConstraintLine(newLine))
	}
}

var condLines = []string{
	`An R1 MUST have an object field named "States"; each field is a "State".`,
	`An R1 which is not an R2 MUST have an object field named "States"; each field is a "State".`,
	`An R1 which is not an R2 or an R3 MUST NOT have a field named "StartAt".`,
	`An R1 which is not an R2, an R3, or an R4 MUST NOT have a field named "StartAt".`,
}

func TestMatcher_CatchConditionalOnConstraint(t *testing.T) {
	t.Parallel()

	excludes := []string{
		``,
		`an R2`,
		`an R2 or an R3`,
		`an R2, an R3, or an R4`,
	}

	m := NewMatcher("R1")
	m.AddRole("R2")
	m.AddRole("R3")
	m.AddRole("R4")

	for i, line := range condLines {
		f := m.BuildConstraint(line)
		val, has := f["excluded"]

		if excludes[i] == "" {
			assert.False(t, has)

			continue
		}

		assert.True(t, has)
		assert.True(t, val == excludes[i])
	}
}

func TestMatcher_MatchReasonablyComplexConstraint(t *testing.T) {
	t.Parallel()

	m := NewMatcher("State")
	s := `A State MUST have a string field named "Type" whose value MUST be one of "Pass", "Succeed", "Fail", "Task", "Choice", "Wait", or "Parallel".`

	assert.True(t, m.IsConstraintLine(s))

	m.AddRole("Retrier")
	s = `A Retrier MAY have a nonnegative-integer field named "MaxAttempts" whose value MUST be less than 99999999.`
	assert.True(t, m.IsConstraintLine(s))
}

func TestMatcher_BuildEnumConstraintObject(t *testing.T) {
	t.Parallel()

	m := NewMatcher("State")
	s := `A State MUST have a string field named "Type" whose value MUST be one of "Pass", "Succeed", "Fail", "Task", "Choice", "Wait", or "Parallel".`

	f := m.BuildConstraint(s)

	assert.True(t, f["role"] == "State")
	assert.True(t, f["modal"] == "MUST")
	assert.True(t, f["type"] == "string")
	assert.True(t, f["field_name"] == "Type")
	_, has := f["relation"]
	assert.False(t, has)
	assert.True(t, f["strings"] == `"Pass", "Succeed", "Fail", "Task", "Choice", "Wait", or "Parallel"`)
	_, has = f["child_type"]
	assert.False(t, has)
}

func TestMatcher_TokenizeStringListsProperly(t *testing.T) {
	t.Parallel()

	m := NewMatcher("x")
	assert.True(t, reflect.DeepEqual(m.TokenizeStrings(`"a"`), []string{`a`}))
	assert.True(t, reflect.DeepEqual(m.TokenizeStrings(`"a" or "b"`), []string{`a`, `b`}))
	assert.True(t, reflect.DeepEqual(m.TokenizeStrings(`"a", "b", or "c"`), []string{`a`, `b`, `c`}))
}

func TestMatcher_BuildARelationalConstraintObject(t *testing.T) {
	t.Parallel()

	m := NewMatcher("Retrier")
	s := `A Retrier MAY have a nonnegative-integer field named "MaxAttempts" whose value MUST be less than 99999999.`
	f := m.BuildConstraint(s)

	assert.True(t, f["role"] == "Retrier")
	assert.True(t, f["modal"] == "MAY")
	assert.True(t, f["type"] == "nonnegative-integer")
	assert.True(t, f["field_name"] == "MaxAttempts")
	_, has := f["strings"]
	assert.False(t, has)
	assert.True(t, f["relation"] == "less than")
	assert.True(t, f["target"] == "99999999")
	_, has = f["child_type"]
	assert.False(t, has)
}

func TestMatcher_BuildConstraintObjectWithChildType(t *testing.T) {
	t.Parallel()

	m := NewMatcher("State Machine")
	s := `A State Machine MUST have an object field named "States"; each field is a "State".`
	f := m.BuildConstraint(s)

	assert.True(t, f["role"] == "State Machine")
	assert.True(t, f["modal"] == "MUST")
	assert.True(t, f["type"] == "object")
	assert.True(t, f["field_name"] == "States")
	assert.True(t, f["child_type"] == "field")
	assert.True(t, f["child_role"] == "State")

	s = `A State Machine MAY have an object field named "Not"; its value is a "FOO".`
	f = m.BuildConstraint(s)

	assert.True(t, f["role"] == "State Machine")
	assert.True(t, f["modal"] == "MAY")
	assert.True(t, f["type"] == "object")
	assert.True(t, f["field_name"] == "Not")
	assert.True(t, f["child_role"] == "FOO")
}
