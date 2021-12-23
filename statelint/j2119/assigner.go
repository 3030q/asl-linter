package j2119

import (
	"regexp"
	"strings"
)

// Assigner looks at the parsed form of the J2119 lines and figures out,
// by looking at which part of the regexes match,
// the assignments of roles to nodes and constraints to roles
type Assigner struct {
	constraints   *RoleConstraints
	roles         *RoleFinder
	matcher       *Matcher
	allowedFields *AllowedFields
}

func NewAssigner(
	constraints *RoleConstraints,
	roles *RoleFinder,
	matcher *Matcher,
	allowedFields *AllowedFields,
) *Assigner {
	return &Assigner{constraints: constraints, roles: roles, matcher: matcher, allowedFields: allowedFields}
}

func (a *Assigner) AssignRoles(assertion map[string]string) {
	if _, ok := assertion["val_match_present"]; ok {
		a.roles.AddFieldValueRole(assertion["role"],
			assertion["fieldtomatch"],
			assertion["valtomatch"],
			assertion["newrole"])
		a.matcher.AddRole(assertion["newrole"])

		return
	}

	if _, ok := assertion["with_a_field"]; ok {
		a.roles.AddFieldPresenceRole(assertion["role"],
			assertion["with_a_field"],
			assertion["newrole"])
		a.matcher.AddRole(assertion["newrole"])

		return
	}

	a.roles.AddIsRole(assertion["role"], assertion["newrole"])
	a.matcher.AddRole(assertion["newrole"])
}

func (a *Assigner) AssignOnlyOneOf(assertion map[string]string) {
	role := assertion["role"]
	values := BreakStringList(assertion["field_list"])
	a.AddConstraint(role, NewOnlyOneConstraint(values), nil)
}

func (a *Assigner) AssignConstraints(assertion map[string]string) {
	role := assertion["role"]
	modal := assertion["modal"]
	typeField, typeOk := assertion["type"]
	fieldName, fieldNameOk := assertion["field_name"]
	fieldListString, fieldListStringOk := assertion["field_list"]
	relation, relationOk := assertion["relation"]
	target := assertion["target"]
	strs, strsOk := assertion["strings"]

	childType, childTypeOk := assertion["child_type"]

	// watch out for conditionals
	var condition *RoleNotPresentedCondition

	if excluded, ok := assertion["excluded"]; ok {
		excludedRoles := BreakRoleList(*a.matcher, excluded)
		condition = NewRoleNotPresentedCondition(excludedRoles)
	}

	if relationOk {
		a.AddRelationConstraint(role, fieldName, relation, target, condition)
	}

	if strsOk {
		// of the form MUST have a <type> field named <field_name> whose value
		// MUST be one of "a", "b", or "c"
		fields := BreakStringList(strs)

		a.AddConstraint(role, NewFieldValueConstraint(fieldName, FieldValueParams{
			Enum: fields,
		}), condition)
	}

	if typeOk {
		a.AddTypeConstraint(role, fieldName, typeField, condition)
	}

	// register allowed fields
	var fieldList []string
	if fieldListStringOk {
		fieldList = BreakStringList(fieldListString)
		for _, field := range fieldList {
			a.allowedFields.SetAllowed(role, field)
		}
	} else if fieldNameOk {
		a.allowedFields.SetAllowed(role, fieldName)
	}

	if modal == "MUST" {
		if fieldListStringOk {
			// Of the form MUST have a <type>? field named one of "a", "b", or "c".
			a.AddConstraint(role, NewHasFieldConstraintFromStringArray(fieldList), condition)
		} else {
			a.AddConstraint(role, NewHasFieldConstraintFromString(fieldName), condition)
		}
	}

	if modal == "MUST NOT" {
		a.AddConstraint(role, NewDoesNotHaveFieldConstraint(fieldName), condition)
	}

	// there can be role defs there too
	if childTypeOk {
		a.matcher.AddRole(assertion["child_role"])

		if childType == "value" {
			a.roles.AddChildRole(role, fieldName, assertion["child_role"])
		} else if childType == "element" || childType == "field" {
			a.roles.AddGrandchildRole(role, fieldName, assertion["child_role"])
		}
	} else {
		anyOrObjectOrArray := (!typeOk) || typeField == "object" || typeField == "array"

		// untyped field without a defined child role
		if fieldNameOk && anyOrObjectOrArray && modal != "MUST NOT" {
			a.roles.AddGrandchildRole(role, fieldName, fieldName)
			a.allowedFields.SetAny(fieldName)
		}
	}
}

func (a *Assigner) AddConstraint(role string, constraint Constrainter, condition *RoleNotPresentedCondition) {
	constraint.AddCondition(condition)
	a.constraints.Add(role, constraint)
}

func (a *Assigner) AddRelationConstraint(
	role string,
	field string,
	relation string,
	target string,
	condition *RoleNotPresentedCondition,
) {
	targetValue := DeduceValue(target)
	params := FieldValueParams{}
	targetValueFloat := GetFloat(targetValue)

	switch relation {
	case "equal to":
		params.IsEqual = true
		params.Equal = targetValueFloat
	case "greater than":
		params.IsFloor = true
		params.Floor = targetValueFloat
	case "less than":
		params.IsCeiling = true
		params.Ceiling = targetValueFloat
	case "greater than or equal to":
		params.IsMin = true
		params.Min = targetValueFloat
	case "less than or equal to":
		params.IsMax = true
		params.Max = targetValueFloat
	}

	a.AddConstraint(role, NewFieldValueConstraint(field, params), condition)
}

var (
	isArrayChecker    = regexp.MustCompile("-array")
	isNullableChecker = regexp.MustCompile("nullable-")
)

func (a *Assigner) AddTypeConstraint(
	role string,
	field string,
	fieldType string,
	condition *RoleNotPresentedCondition,
) {
	isArray := isArrayChecker.MatchString(fieldType)
	isNullable := isNullableChecker.MatchString(fieldType)

	for _, part := range strings.Split(fieldType, "-") {
		switch part {
		case "object":
			a.AddConstraint(role, NewFieldTypeConstraint(field, Object, isArray, isNullable), condition)
		case "string":
			a.AddConstraint(role, NewFieldTypeConstraint(field, String, isArray, isNullable), condition)
		case "URI":
			a.AddConstraint(role, NewFieldTypeConstraint(field, URI, isArray, isNullable), condition)
		case "boolean":
			a.AddConstraint(role, NewFieldTypeConstraint(field, Bool, isArray, isNullable), condition)
		case "numeric":
			a.AddConstraint(role, NewFieldTypeConstraint(field, Numeric, isArray, isNullable), condition)
		case "integer":
			a.AddConstraint(role, NewFieldTypeConstraint(field, Integer, isArray, isNullable), condition)
		case "float":
			a.AddConstraint(role, NewFieldTypeConstraint(field, Float, isArray, isNullable), condition)
		case "timestamp":
			a.AddConstraint(role, NewFieldTypeConstraint(field, Timestamp, isArray, isNullable), condition)
		case "JSONPath":
			a.AddConstraint(role, NewFieldTypeConstraint(field, JSONPath, isArray, isNullable), condition)
		case "referencePath":
			a.AddConstraint(
				role,
				NewFieldTypeConstraint(field, ReferencePath, isArray, isNullable),
				condition,
			)
		case "positive":
			a.AddConstraint(role, NewFieldValueConstraint(field, FieldValueParams{
				IsFloor: true,
				Floor:   0,
			}), condition)
		case "nonnegative":
			a.AddConstraint(role, NewFieldValueConstraint(field, FieldValueParams{
				IsMin: true,
				Min:   0,
			}), condition)
		case "negative":
			a.AddConstraint(role, NewFieldValueConstraint(field, FieldValueParams{
				IsCeiling: true,
				Ceiling:   0,
			}), condition)
		case "nonempty":
			a.AddConstraint(role, NewNonEmptyConstraint(field), condition)
		}
	}
}
