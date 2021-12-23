package j2119

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssigner_AttachConditionToConstraint(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":       "R",
		"modal":      "MUST",
		"field_name": "foo",
		"excluded":   "an A, a B, or a C",
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()

	a := NewAssigner(constraints, rf, matcher, allowedFields)

	for _, x := range []string{"A", "B", "C"} {
		matcher.AddRole(x)
	}

	a.AssignConstraints(assertion)

	retrieved := constraints.Get("R")
	c := retrieved[0]
	node := NewNodeCreateHelper(t, `{"a":1}`)

	assert.IsType(t, &HasFieldConstraint{}, c)

	for _, role := range []string{"A", "B", "C"} {
		assert.False(t, c.Applies(node, []string{role}))
	}

	assert.True(t, c.Applies(node, []string{"foo"}))
}

func TestAssigner_HandleNonZero___LessThanConstraintProperly(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":       "R",
		"modal":      "May",
		"type":       "nonnegative-integer",
		"field_name": "MaxAttempts",
		"relation":   "less than",
		"target":     "99999999",
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()
	a := NewAssigner(constraints, rf, matcher, allowedFields)

	a.AssignConstraints(assertion)

	retrieved := constraints.Get("R")

	assert.Equal(t, 3, len(retrieved))
}

func TestAssigner_AssignOnlyOneOfConstraintProperly(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":       "R",
		"field_list": `"foo", "bar", and "baz"`,
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()
	a := NewAssigner(constraints, rf, matcher, allowedFields)

	a.AssignOnlyOneOf(assertion)

	retrieved := constraints.Get("R")

	assert.IsType(t, &OnlyOneConstraint{}, retrieved[0])
}

func TestAddHasFieldConstraintIfThereMUST(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":       "R",
		"modal":      "MUST",
		"field_name": `foo`,
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()
	a := NewAssigner(constraints, rf, matcher, allowedFields)

	a.AssignConstraints(assertion)

	retrieved := constraints.Get("R")

	assert.IsType(t, &HasFieldConstraint{}, retrieved[0])
}

func TestAssigner_AddDoesNotHaveFieldConstraintIfMUSTNOT(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":       "R",
		"modal":      "MUST NOT",
		"field_name": `foo`,
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()
	a := NewAssigner(constraints, rf, matcher, allowedFields)

	a.AssignConstraints(assertion)

	retrieved := constraints.Get("R")

	assert.IsType(t, &DoesNotHaveFieldConstraint{}, retrieved[0])
}

func TestAssigner_ManageComplexTypeConstraint(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":       "R",
		"modal":      "MUST",
		"field_name": `foo`,
		"type":       "nonnegative-float",
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()
	a := NewAssigner(constraints, rf, matcher, allowedFields)

	a.AssignConstraints(assertion)

	retrieved := constraints.Get("R")

	// add all added constraints to arrays by type
	var hasFieldConstraints []*HasFieldConstraint
	var fieldTypeConstraint []*FieldTypeConstraint
	var fieldValueConstraint []*FieldValueConstraint

	for _, constrainter := range retrieved {
		switch reflect.TypeOf(constrainter) {
		case reflect.TypeOf(&HasFieldConstraint{}):
			hasFieldConstraints = append(hasFieldConstraints, constrainter.(*HasFieldConstraint))
		case reflect.TypeOf(&FieldTypeConstraint{}):
			fieldTypeConstraint = append(fieldTypeConstraint, constrainter.(*FieldTypeConstraint))
		case reflect.TypeOf(&FieldValueConstraint{}):
			fieldValueConstraint = append(fieldValueConstraint, constrainter.(*FieldValueConstraint))
		}
	}

	assert.Equal(t, 1, len(hasFieldConstraints))
	assert.Equal(t, 1, len(fieldTypeConstraint))
	assert.Equal(t, 1, len(fieldValueConstraint))
}

func TestAssigner_RecordRelationalConstraint(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":       "R",
		"modal":      "MUST",
		"field_name": `foo`,
		"type":       "nonnegative-float",
		"relation":   "less than",
		"target":     "123.09",
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()
	a := NewAssigner(constraints, rf, matcher, allowedFields)

	a.AssignConstraints(assertion)

	retrieved := constraints.Get("R")

	// add all added constraints to arrays by type
	var hasFieldConstraints []*HasFieldConstraint
	var fieldTypeConstraint []*FieldTypeConstraint
	var fieldValueConstraint []*FieldValueConstraint

	for _, constrainter := range retrieved {
		switch reflect.TypeOf(constrainter) {
		case reflect.TypeOf(&HasFieldConstraint{}):
			hasFieldConstraints = append(hasFieldConstraints, constrainter.(*HasFieldConstraint))
		case reflect.TypeOf(&FieldTypeConstraint{}):
			fieldTypeConstraint = append(fieldTypeConstraint, constrainter.(*FieldTypeConstraint))
		case reflect.TypeOf(&FieldValueConstraint{}):
			fieldValueConstraint = append(fieldValueConstraint, constrainter.(*FieldValueConstraint))
		}
	}

	assert.Equal(t, 1, len(hasFieldConstraints))
	assert.Equal(t, 1, len(fieldTypeConstraint))
	assert.Equal(t, 2, len(fieldValueConstraint))
}

func TestAssigner_RecordIsRole(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":    "R",
		"newrole": "S",
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()
	a := NewAssigner(constraints, rf, matcher, allowedFields)

	a.AssignRoles(assertion)

	node := NewNodeCreateHelper(t, `{"a": 3}`)
	roles := []string{"R"}

	moreRoles := rf.FindMoreRoles(node, roles)

	assert.Equal(t, 2, len(moreRoles))

	hasSRole := false

	for _, role := range moreRoles {
		if role == "S" {
			hasSRole = true

			break
		}
	}

	assert.True(t, hasSRole)
}

func TestAssigner_CorrectlyAssignFieldValueRole(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":              "R",
		"fieldtomatch":      "f1",
		"valtomatch":        "33",
		"newrole":           "S",
		"val_match_present": "true",
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()
	a := NewAssigner(constraints, rf, matcher, allowedFields)

	a.AssignRoles(assertion)

	node := NewNodeCreateHelper(t, `{"f1": 33}`)
	roles := []string{"R"}
	moreRoles := rf.FindMoreRoles(node, roles)

	hasSRole := false

	for _, role := range moreRoles {
		if role == "S" {
			hasSRole = true

			break
		}
	}

	assert.True(t, hasSRole)
}

func TestAssigner_ProcessChildRoleInAssertion(t *testing.T) {
	t.Parallel()

	assertion := map[string]string{
		"role":       "R",
		"modal":      "MUST",
		"field_name": "a",
		"child_type": "field",
		"child_role": "bar",
	}
	constraints := NewRoleConstraints()
	rf := NewRoleFinder()
	matcher := NewMatcher("x")
	allowedFields := NewAllowedFields()
	a := NewAssigner(constraints, rf, matcher, allowedFields)

	a.AssignConstraints(assertion)

	roles := []string{"R"}
	fieldRoles := rf.FindGrandchildRoles(roles, "a")

	assert.Equal(t, 1, len(fieldRoles))
	assert.Equal(t, "bar", fieldRoles[0])
}
