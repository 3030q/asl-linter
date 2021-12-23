package j2119

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeParser struct {
	rc *RoleConstraints
	rf *RoleFinder
}

func (f FakeParser) FindMoreRoles(node Node, roles []string) []string {
	return f.rf.FindMoreRoles(node, roles)
}

func (f FakeParser) FindGrandchildRoles(roles []string, name string) []string {
	return f.rf.FindGrandchildRoles(roles, name)
}

func (f FakeParser) FindChildRoles(roles []string, name string) []string {
	return f.rf.FindChildRoles(roles, name)
}

func (f FakeParser) GetConstraints(role string) []Constrainter {
	return f.rc.Get(role)
}

func (f FakeParser) IsFieldAllowed(_ []string, _ string) bool {
	return true
}

func (f FakeParser) IsAllowsAny(_ []string) bool {
	return false
}

func NewFakeParser(rc *RoleConstraints, rf *RoleFinder) *FakeParser {
	return &FakeParser{rc: rc, rf: rf}
}

func TestNodeValidator_ReportProblemsWithFaultyFields(t *testing.T) {
	t.Parallel()

	rf := NewRoleFinder()
	rc := NewRoleConstraints()
	nv := NewNodeValidator(NewFakeParser(rc, rf))

	roles := []string{"Role1"}

	// among fields
	// 'a' should exist
	// 'b' should not exist
	// 'c' should be a float
	// 'd' should be an integer
	// 'e' should be a number
	// 'f' should be between 0 and 5
	node := NewNodeCreateHelper(t, `{"b":1,"c":"float","d":0.3,"e":true,"f":10}`)

	constraints := []Constrainter{
		NewHasFieldConstraintFromString("a"),
		NewDoesNotHaveFieldConstraint("b"),
		NewFieldTypeConstraint("c", Float, false, false),
		NewFieldTypeConstraint("d", Integer, false, false),
		NewFieldTypeConstraint("e", Numeric, false, false),
		NewFieldValueConstraint("f", FieldValueParams{
			IsMin: true,
			IsMax: true,
			Min:   0,
			Max:   5,
		}),
	}

	for _, constraint := range constraints {
		rc.Add("Role1", constraint)
	}

	problems := NewProblems()
	nv.Validate(node, "x.y", roles, problems)

	assert.Equal(t, len(constraints), problems.Len())
}
