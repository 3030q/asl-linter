package j2119

import "reflect"

// RoleFinder This is about figuring out which roles apply to a node and
// potentially to its children in object and array valued fields
type RoleFinder struct {
	// roles of the form: A Foo is a Bar
	IsRoles map[string][]string

	// roles of the form: If an object with role X has a field Y, then
	// it has role R
	// map[role][field_name] => child_role
	FieldPresentRoles map[string]map[string]string

	// roles of the form: If an object with role X has a field Y with
	// value Z, it has role R
	// map[role][field_name][field_val] => child_role
	FieldValueRoles map[string]map[string]map[interface{}]string

	// roles of the form: If an object with role X has field Y which
	// is an object/array, the object-files/array-elements have role R
	GrandchildRoles map[string]map[string]string

	// roles of the form: If an object with role X has field Y which
	// is an object, that object has role R
	ChildRoles map[string]map[string]string
}

func NewRoleFinder() *RoleFinder {
	return &RoleFinder{
		IsRoles:           make(map[string][]string),
		FieldPresentRoles: make(map[string]map[string]string),
		FieldValueRoles:   make(map[string]map[string]map[interface{}]string),
		GrandchildRoles:   make(map[string]map[string]string),
		ChildRoles:        make(map[string]map[string]string),
	}
}

func (r *RoleFinder) AddIsRole(role string, otherRole string) {
	r.IsRoles[role] = append(r.IsRoles[role], otherRole)
}

func (r *RoleFinder) AddFieldValueRole(role string, fieldName string, fieldValue string, newRole string) {
	if r.FieldValueRoles[role] == nil {
		r.FieldValueRoles[role] = make(map[string]map[interface{}]string)
	}

	if r.FieldValueRoles[role][fieldName] == nil {
		r.FieldValueRoles[role][fieldName] = make(map[interface{}]string)
	}

	value := DeduceValue(fieldValue)
	r.FieldValueRoles[role][fieldName][value] = newRole
}

func (r *RoleFinder) AddFieldPresenceRole(role string, fieldName string, childRole string) {
	if r.FieldPresentRoles[role] == nil {
		r.FieldPresentRoles[role] = make(map[string]string)
	}

	r.FieldPresentRoles[role][fieldName] = childRole
}

func (r *RoleFinder) AddChildRole(role string, fieldName string, childRole string) {
	if r.ChildRoles[role] == nil {
		r.ChildRoles[role] = make(map[string]string)
	}

	r.ChildRoles[role][fieldName] = childRole
}

func (r *RoleFinder) AddGrandchildRole(role string, fieldName string, childRole string) {
	if r.GrandchildRoles[role] == nil {
		r.GrandchildRoles[role] = make(map[string]string)
	}

	r.GrandchildRoles[role][fieldName] = childRole
}

// FindMoreRoles Consider a node which has one or more roles. It may have more
// roles based on the presence or value of child nodes. This method
// adds any such roles to the "roles" list
func (r *RoleFinder) FindMoreRoles(node Node, roles []string) []string {
	var result []string
	result = append(result, roles...)

	// find roles depending on field values
	for _, role := range result {
		perFieldName, ok := r.FieldValueRoles[role]
		if ok {
			for fieldName, valueRoles := range perFieldName {
				for fieldValue, childRole := range valueRoles {
					if !node.HasNode(fieldName) {
						break
					}

					v := node.GetNode(fieldName).Value()
					if reflect.DeepEqual(fieldValue, v) {
						result = append(result, childRole)
					}
				}
			}
		}
	}

	// find roles depending on field presence
	for _, role := range result {
		perFieldName, ok := r.FieldPresentRoles[role]
		if ok {
			for fieldName, childRole := range perFieldName {
				if node.HasNode(fieldName) {
					result = append(result, childRole)
				}
			}
		}
	}

	// is roles
	for _, role := range result {
		otherRoles, ok := r.IsRoles[role]
		if ok {
			result = append(result, otherRoles...)
		}
	}

	return result
}

// FindChildRoles A node has a role, and one of its fields might be object-valued
// and that value is given a role
func (r *RoleFinder) FindChildRoles(roles []string, fieldName string) []string {
	var result []string

	for _, role := range roles {
		if _, ok := r.ChildRoles[role]; ok {
			if value, ok := r.ChildRoles[role][fieldName]; ok {
				result = append(result, value)
			}
		}
	}

	return result
}

// FindGrandchildRoles A node has a role, and one of its field is an object or an
// array whose fields or elements are given a role
func (r *RoleFinder) FindGrandchildRoles(roles []string, fieldName string) []string {
	var result []string

	for _, role := range roles {
		if _, ok := r.GrandchildRoles[role]; ok {
			if value, ok := r.GrandchildRoles[role][fieldName]; ok {
				result = append(result, value)
			}
		}
	}

	return result
}
