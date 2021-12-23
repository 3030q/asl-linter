package j2119

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleConstraints_SuccessfullyRememberConstraints(t *testing.T) {
	t.Parallel()

	m := NewRoleConstraints()
	c1 := NewHasFieldConstraintFromString("foo")
	c2 := NewDoesNotHaveFieldConstraint("bar")

	m.Add("MyRole", c1)
	m.Add("MyRole", c2)
	m.Add("OtherRole", c1)

	r := m.Get("MyRole")
	for _, constrainter := range r {
		if !(constrainter == c1 || constrainter == c2) {
			assert.Fail(t, "bad constraint")
		}
	}

	assert.Equal(t, 2, len(r))

	r = m.Get("OtherRole")
	assert.Equal(t, 1, len(r))
	assert.Equal(t, c1, r[0])

	assert.Equal(t, 0, len(m.Get("No Constraint")))
}
