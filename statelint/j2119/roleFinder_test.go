package j2119

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	defaultRole      = "MyRole"
	defaultAddedRole = "NewRole"
)

func TestRoleFinder_SuccessfullyAssignAdditionalRoleBasedOnRole(t *testing.T) {
	t.Parallel()

	r := NewRoleFinder()
	node := NewNodeCreateHelper(t, `{"a": 3}`)
	roles := []string{defaultRole}

	r.AddIsRole(defaultRole, defaultAddedRole)

	moreRoles := r.FindMoreRoles(node, roles)

	assert.Equal(t, 2, len(moreRoles))

	hasRole := false

	for _, role := range moreRoles {
		if role == defaultAddedRole {
			hasRole = true

			break
		}
	}

	assert.True(t, hasRole)
}

func TestRoleFinder_SuccessfullyAssignRoleBasedOnFieldValue(t *testing.T) {
	t.Parallel()

	r := NewRoleFinder()
	node := NewNodeCreateHelper(t, `{"a": 3}`)

	r.AddFieldValueRole(defaultRole, "a", "3", defaultAddedRole)

	roles := []string{defaultRole}

	moreRoles := r.FindMoreRoles(node, roles)

	assert.Equal(t, 2, len(moreRoles))

	hasRole := false

	for _, role := range moreRoles {
		if role == defaultAddedRole {
			hasRole = true

			break
		}
	}

	assert.True(t, hasRole)
}

func TestRoleFinder_SuccessfullyAssignRoleBasedOnFieldPresence(t *testing.T) {
	t.Parallel()

	r := NewRoleFinder()
	node := NewNodeCreateHelper(t, `{"a": 3}`)

	r.AddFieldPresenceRole(defaultRole, "a", defaultAddedRole)

	roles := []string{defaultRole}

	moreRoles := r.FindMoreRoles(node, roles)

	assert.Equal(t, 2, len(moreRoles))

	hasRole := false

	for _, role := range moreRoles {
		if role == defaultAddedRole {
			hasRole = true

			break
		}
	}

	assert.True(t, hasRole)
}

func TestRoleFinder_SuccessfullyAddRoleToGrandchildrenFieldBasedOnName(t *testing.T) {
	t.Parallel()

	r := NewRoleFinder()
	r.AddChildRole(defaultRole, "a", defaultAddedRole)

	roles := []string{defaultRole}

	childRoles := r.FindChildRoles(roles, "a")

	assert.Equal(t, 1, len(childRoles))

	hasRole := false

	for _, role := range childRoles {
		if role == defaultAddedRole {
			hasRole = true

			break
		}
	}

	assert.True(t, hasRole)
}

func TestRoleFinder_SuccessfullyAddRoleToChildFieldBasedOnName(t *testing.T) {
	t.Parallel()

	r := NewRoleFinder()
	r.AddGrandchildRole(defaultRole, "a", defaultAddedRole)

	roles := []string{defaultRole}

	grandchildRoles := r.FindGrandchildRoles(roles, "a")

	assert.Equal(t, 1, len(grandchildRoles))

	hasRole := false

	for _, role := range grandchildRoles {
		if role == defaultAddedRole {
			hasRole = true

			break
		}
	}

	assert.True(t, hasRole)
}
