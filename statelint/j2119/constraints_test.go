package j2119

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewNodeCreateHelper(t *testing.T, jsonString string) Node {
	t.Helper()

	var j interface{}
	err := json.Unmarshal([]byte(jsonString), &j)

	if err != nil {
		t.Fatal("cant unmarshal json")
	}

	return *NewNode(j)
}

func TestConstraints_LoadAndEvaluateCondition(t *testing.T) {
	t.Parallel()

	c := NewHasFieldConstraintFromString("foo")

	node := NewNodeCreateHelper(t, `{ "bar": 1 }`)

	assert.True(t, c.Applies(node, []string{"foo"}))

	cond := NewRoleNotPresentedCondition([]string{"foo", "bar"})
	c.AddCondition(cond)

	assert.False(t, c.Applies(node, []string{"foo"}))
	assert.True(t, c.Applies(node, []string{"baz"}))
}

func TestHasFieldConstraint_SuccessfullyDetectMissingField(t *testing.T) {
	t.Parallel()

	c := NewHasFieldConstraintFromString("foo")

	node := NewNodeCreateHelper(t, `{ "bar": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestHasFieldConstraint_AcceptNodeWithRequiredFieldPresent(t *testing.T) {
	t.Parallel()

	c := NewHasFieldConstraintFromString("bar")
	node := NewNodeCreateHelper(t, `{ "bar": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestNonEmptyConstraint_BypassAbsentField(t *testing.T) {
	t.Parallel()

	c := NewNonEmptyConstraint("foo")
	node := NewNodeCreateHelper(t, `{ "bar": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestNonEmptyConstraint_BypassNonArrayField(t *testing.T) {
	t.Parallel()

	c := NewNonEmptyConstraint("foo")
	node := NewNodeCreateHelper(t, `{ "foo": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestNonEmptyConstraint_OkNonEmptyArray(t *testing.T) {
	t.Parallel()

	c := NewNonEmptyConstraint("foo")
	node := NewNodeCreateHelper(t, `{ "foo": [ 1 ] }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestNonEmptyConstraint_CatchEmptyArray(t *testing.T) {
	t.Parallel()

	c := NewNonEmptyConstraint("foo")
	node := NewNodeCreateHelper(t, `{ "foo": [ ] }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestDoesNotHaveFieldConstraint_SuccessfullyDetectForbiddenField(t *testing.T) {
	t.Parallel()

	c := NewDoesNotHaveFieldConstraint("foo")
	node := NewNodeCreateHelper(t, `{ "foo": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestDoesNotHaveFieldConstraint_AcceptNodeWithRequiredFieldPresent(t *testing.T) {
	t.Parallel()

	c := NewDoesNotHaveFieldConstraint("bar")
	node := NewNodeCreateHelper(t, `{ "foo": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestFieldValueConstraint_SilentNoOpExitIfTheFieldIsntThere(t *testing.T) {
	t.Parallel()

	c := NewFieldValueConstraint("bar", FieldValueParams{})
	node := NewNodeCreateHelper(t, `{ "foo": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestFieldValueConstraint_DetectViolationEnumPolicy(t *testing.T) {
	t.Parallel()

	c := NewFieldValueConstraint("foo", FieldValueParams{
		Enum: []string{"1", "2", "3"},
	})
	node := NewNodeCreateHelper(t, `{ "foo": 5 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestFieldValueConstraint_DetectBrokenEquals(t *testing.T) {
	t.Parallel()

	c := NewFieldValueConstraint("foo", FieldValueParams{
		IsEqual: true,
		Equal:   12,
	})
	node := NewNodeCreateHelper(t, `{ "foo": 12 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())

	node = NewNodeCreateHelper(t, `{ "foo": 3 }`)

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestFieldValueConstraint_DoMinRight(t *testing.T) {
	t.Parallel()

	c := NewFieldValueConstraint("foo", FieldValueParams{
		IsMin: true,
		Min:   1,
	})
	node := NewNodeCreateHelper(t, `{ "foo": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())

	node = NewNodeCreateHelper(t, `{ "foo": 0 }`)

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestFieldValueConstraint_DetectBrokenFloor(t *testing.T) {
	t.Parallel()

	c := NewFieldValueConstraint("foo", FieldValueParams{
		IsFloor: true,
		Floor:   1,
	})
	node := NewNodeCreateHelper(t, `{ "foo": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestFieldValueConstraint_DetectBrokenCeiling(t *testing.T) {
	t.Parallel()

	c := NewFieldValueConstraint("foo", FieldValueParams{
		IsCeiling: true,
		Ceiling:   3,
	})
	node := NewNodeCreateHelper(t, `{ "foo": 3 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestFieldValueConstraint_DoMaxRight(t *testing.T) {
	t.Parallel()

	c := NewFieldValueConstraint("foo", FieldValueParams{
		IsMax: true,
		Max:   3,
	})
	node := NewNodeCreateHelper(t, `{ "foo": 3 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())

	node = NewNodeCreateHelper(t, `{ "foo": 4 }`)

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestFieldValueConstraint_AcceptMinMax(t *testing.T) {
	t.Parallel()

	c := NewFieldValueConstraint("foo", FieldValueParams{
		IsMin: true,
		IsMax: true,
		Min:   0,
		Max:   3,
	})
	node := NewNodeCreateHelper(t, `{ "foo": 1 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestOnlyOneOfConstraint_DetectMoreThanOneErrors(t *testing.T) {
	t.Parallel()

	c := NewOnlyOneConstraint([]string{"foo", "bar", "baz"})
	node := NewNodeCreateHelper(t, `{ "foo": 1, "bar": 2 }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestFieldTypeConstraint_SilentNoOpExitIfFieldIsntThere(t *testing.T) {
	t.Parallel()

	c := NewFieldTypeConstraint("foo", Integer, false, false)
	node := NewNodeCreateHelper(t, `{ "bar": 1}`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestFieldTypeConstraint_ApproveCorrectTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		valueType ValueType
		value     interface{}
	}{
		{String, `"foo"`},
		{Integer, "3"},
		{Float, "0.33"},
		{Bool, "false"},
		{Timestamp, `"2016-03-14T01:59:00Z"`},
		{Object, `{ "a": 1 }`},
		{Array, `[ 3, 4 ]`},
		{JSONPath, `"$.a.c[2,3]"`},
		{ReferencePath, `"$.a['b'].d[3]"`},
	}

	for i, testCase := range testCases {
		c := NewFieldTypeConstraint("foo", testCase.valueType, false, false)
		node := NewNodeCreateHelper(t, fmt.Sprintf(`{ "foo": %s}`, testCase.value))
		problems := NewProblems()

		c.Check(node, "a.b.c", problems)

		assert.Equal(t, 0, problems.Len(), fmt.Sprintf("test case number %d", i))
	}
}

func TestFieldTypeConstraint_FindIncorrectTypesInAnArrayField(t *testing.T) {
	t.Parallel()

	c := NewFieldTypeConstraint("a", Integer, false, false)
	node := NewNodeCreateHelper(t, `{ "a": [ 1, 2, "foo", 4 ] }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestFieldTypeConstraint_FlagIncorrectTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		valueType ValueType
		value     interface{}
	}{
		{String, `33`},
		{Integer, `"foo"`},
		{Float, `true`},
		{Bool, "null"},
		{Timestamp, `"2x16-03-14T01:59:00Z"`},
		{JSONPath, `"blibble"`},
		{ReferencePath, `"$.a.*"`},
	}

	for i, testCase := range testCases {
		c := NewFieldTypeConstraint("foo", testCase.valueType, false, false)
		node := NewNodeCreateHelper(t, fmt.Sprintf(`{ "foo": %s}`, testCase.value))
		problems := NewProblems()

		c.Check(node, "a.b.c", problems)

		assert.Equal(t, 1, problems.Len(), fmt.Sprintf("test case number %d", i))
	}
}

func TestFieldTypeConstraint_HandleNullableCorrectly(t *testing.T) {
	t.Parallel()

	c := NewFieldTypeConstraint("a", String, false, false)
	node := NewNodeCreateHelper(t, `{ "a": null }`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())

	c = NewFieldTypeConstraint("a", String, false, true)
	problems = NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestFieldTypeConstraint_HandleArrayNestingConstraints(t *testing.T) {
	t.Parallel()

	c := NewFieldTypeConstraint("foo", Array, false, false)
	node := NewNodeCreateHelper(t, `{"foo": 1}`)
	problems := NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())

	c = NewFieldTypeConstraint("foo", Integer, true, true)
	node = NewNodeCreateHelper(t, `{"foo": [ "bar" ] }`)
	problems = NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 1, problems.Len())

	node = NewNodeCreateHelper(t, `{"foo": [ 1 ] }`)
	problems = NewProblems()

	c.Check(node, "a.b.c", problems)

	assert.Equal(t, 0, problems.Len())
}
