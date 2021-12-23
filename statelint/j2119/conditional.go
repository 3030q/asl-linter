package j2119

import "fmt"

type RoleNotPresentedCondition struct {
	excludeRoles []string
}

func NewRoleNotPresentedCondition(excludeRoles []string) *RoleNotPresentedCondition {
	return &RoleNotPresentedCondition{excludeRoles: excludeRoles}
}

func (r *RoleNotPresentedCondition) String() string {
	return fmt.Sprintf("excluded roles: %v", r.excludeRoles)
}

func (r *RoleNotPresentedCondition) IsConstraintApplies(_ Node, roles []string) bool {
	for _, exRole := range r.excludeRoles {
		for _, role := range roles {
			if exRole == role {
				return false
			}
		}
	}

	return true
}

func IsAppliesAny(conditions []RoleNotPresentedCondition, node Node, roles []string) bool {
	for _, condition := range conditions {
		if condition.IsConstraintApplies(node, roles) {
			return true
		}
	}

	return false
}
