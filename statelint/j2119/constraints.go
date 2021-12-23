package j2119

import (
	"fmt"
	"reflect"
	"regexp"
	"statelint/localization"
)

type Constrainter interface {
	Check(node Node, path string, problems *Problems)
	AddCondition(constraint *RoleNotPresentedCondition)
	Applies(node Node, roles []string) bool
}

type Constraint struct {
	Conditions []RoleNotPresentedCondition
}

func NewConstraint() *Constraint {
	return &Constraint{
		Conditions: make([]RoleNotPresentedCondition, 0),
	}
}

func (c *Constraint) AddCondition(constraint *RoleNotPresentedCondition) {
	if constraint != nil {
		c.Conditions = append(c.Conditions, *constraint)
	}
}

func (c *Constraint) Applies(node Node, roles []string) bool {
	return len(c.Conditions) == 0 || IsAppliesAny(c.Conditions, node, roles)
}

type OnlyOneConstraint struct {
	Constraint
	fields []string
}

func NewOnlyOneConstraint(fields []string) *OnlyOneConstraint {
	return &OnlyOneConstraint{
		Constraint: *NewConstraint(),
		fields:     fields,
	}
}

func (o *OnlyOneConstraint) Check(node Node, path string, problems *Problems) {
	join := make([]string, 0)
	keys := node.Keys()

	for _, field := range o.fields {
		for _, key := range keys {
			if field == key {
				join = append(join, field)
			}
		}
	}

	if len(join) > 1 {
		problemStr := localization.GetLocalizerOrPanic().GetString("OnlyOneConstraint")
		problems.Append(fmt.Sprintf(problemStr, path, o.fields))
	}
}

type NonEmptyConstraint struct {
	Constraint
	name string
}

func NewNonEmptyConstraint(name string) *NonEmptyConstraint {
	return &NonEmptyConstraint{
		Constraint: *NewConstraint(),
		name:       name,
	}
}

func (n *NonEmptyConstraint) String() string {
	var conds string
	if len(n.Conditions) != 0 {
		conds = fmt.Sprintf("%d conditions", len(n.Conditions))
	}

	return fmt.Sprintf("<Array field %s should not be empty %s>", n.name, conds)
}

func (n *NonEmptyConstraint) Check(node Node, path string, problems *Problems) {
	if node.HasNode(n.name) &&
		node.GetNode(n.name).Is(Array) &&
		len(node.GetNode(n.name).ValueToArray()) == 0 {
		problemStr := localization.GetLocalizerOrPanic().GetString("NonEmptyConstraint")
		problems.Append(fmt.Sprintf(problemStr, path, n.name))
	}
}

type HasFieldConstraint struct {
	Constraint
	names []string
}

func NewHasFieldConstraintFromString(name string) *HasFieldConstraint {
	return &HasFieldConstraint{
		Constraint: *NewConstraint(),
		names:      []string{name},
	}
}

func NewHasFieldConstraintFromStringArray(names []string) *HasFieldConstraint {
	return &HasFieldConstraint{
		Constraint: *NewConstraint(),
		names:      names,
	}
}

func (h *HasFieldConstraint) String() string {
	var conds string
	if len(h.Conditions) != 0 {
		conds = fmt.Sprintf("%d conditions", len(h.Conditions))
	}

	return fmt.Sprintf("<Field %v should not be empty %s>", h.names, conds)
}

func (h *HasFieldConstraint) Check(node Node, path string, problems *Problems) {
	if !node.Is(Object) {
		return
	}

	join := make([]string, 0)
	keys := node.Keys()

	for _, field := range h.names {
		for _, key := range keys {
			if field == key {
				join = append(join, field)
			}
		}
	}

	if len(join) == 0 {
		if len(h.names) == 1 {
			problemStr := localization.GetLocalizerOrPanic().GetString("HasFieldConstraintSingle")
			problems.Append(fmt.Sprintf(problemStr, path, h.names[0]))
		} else {
			problemStr := localization.GetLocalizerOrPanic().GetString("HasFieldConstraintMultiple")
			problems.Append(fmt.Sprintf(problemStr, path, h.names))
		}
	}
}

type DoesNotHaveFieldConstraint struct {
	Constraint
	name string
}

func NewDoesNotHaveFieldConstraint(name string) *DoesNotHaveFieldConstraint {
	return &DoesNotHaveFieldConstraint{
		Constraint: *NewConstraint(),
		name:       name,
	}
}

func (d *DoesNotHaveFieldConstraint) String() string {
	var conds string
	if len(d.Conditions) != 0 {
		conds = fmt.Sprintf("%d conditions", len(d.Conditions))
	}

	return fmt.Sprintf("<Field %s should be absent %s>", d.name, conds)
}

func (d *DoesNotHaveFieldConstraint) Check(node Node, path string, problems *Problems) {
	if node.HasNode(d.name) {
		problemStr := localization.GetLocalizerOrPanic().GetString("DoesNotHaveFieldConstraint")
		problems.Append(fmt.Sprintf(problemStr, path, d.name))
	}
}

type FieldTypeConstraint struct {
	Constraint
	name       string
	fieldType  ValueType
	isArray    bool
	isNullable bool
}

func NewFieldTypeConstraint(name string, fieldType ValueType, isArray bool, isNullable bool) *FieldTypeConstraint {
	return &FieldTypeConstraint{
		Constraint: *NewConstraint(),
		name:       name,
		fieldType:  fieldType,
		isArray:    isArray,
		isNullable: isNullable,
	}
}

func (f *FieldTypeConstraint) Check(node Node, path string, problems *Problems) {
	if !node.HasNode(f.name) {
		return
	}

	path = fmt.Sprintf("%s.%s", path, f.name)
	valueNode := node.GetNode(f.name)

	if valueNode.IsNull() {
		if !f.isNullable {
			problemStr := localization.GetLocalizerOrPanic().GetString("FieldTypeConstraint")
			problems.Append(fmt.Sprintf(problemStr, path))
		}

		return
	}

	if f.isArray {
		if valueNode.Is(Array) {
			arr := valueNode.ValueToArray()
			for i, item := range arr {
				f.ValueCheck(item, fmt.Sprintf("%s[%d]", path, i), problems)
			}

			return
		}

		f.ReportValue(path, node, "an Array", problems)

		return
	}

	f.ValueCheck(*valueNode, path, problems)
}

var uriMatcher = regexp.MustCompile(`^[a-z]+:`)

func (f *FieldTypeConstraint) ValueCheck(value Node, path string, problems *Problems) {
	switch f.fieldType {
	case Object:
		if value.Is(Object) {
			return
		}
	case Array:
		if value.Is(Array) {
			return
		}
	case String:
		if value.Is(String) {
			return
		}
	case Integer:
		if value.Is(Integer) {
			return
		}
	case Float:
		if value.Is(Float) {
			return
		}
	case Bool:
		if value.Is(Bool) {
			return
		}
	case Numeric:
		if value.Is(Numeric) {
			return
		}
	case JSONPath:
		if value.Is(JSONPath) {
			return
		}
	case ReferencePath:
		if value.Is(ReferencePath) {
			return
		}
	case Timestamp:
		if value.Is(Timestamp) {
			return
		}
	case URI:
		if value.Is(URI) {
			return
		}
	default:
		panic(fmt.Sprintf("unrecognized field type %s", f.fieldType))
	}

	f.ReportValue(path, value, f.fieldType, problems)
}

func (f *FieldTypeConstraint) ReportValue(path string, value Node, message ValueType, problems *Problems) {
	problemStr := localization.GetLocalizerOrPanic().GetString("FieldTypeConstraintReport")
	problems.Append(fmt.Sprintf(problemStr, path, value.Types(), message))
}

type FieldValueParams struct {
	Enum []string

	IsEqual   bool
	IsFloor   bool
	IsMin     bool
	IsCeiling bool
	IsMax     bool

	Equal   float64
	Floor   float64
	Min     float64
	Ceiling float64
	Max     float64
}

type FieldValueConstraint struct {
	Constraint
	name   string
	params FieldValueParams
}

func NewFieldValueConstraint(name string, params FieldValueParams) *FieldValueConstraint {
	return &FieldValueConstraint{
		Constraint: *NewConstraint(),
		name:       name,
		params:     params,
	}
}

func (f *FieldValueConstraint) String() string {
	var conds string
	if len(f.Conditions) != 0 {
		conds = fmt.Sprintf("%d conditions", len(f.Conditions))
	}

	return fmt.Sprintf("<Field %s has constraints %v%s>", f.name, f.params, conds)
}

func (f *FieldValueConstraint) Check(node Node, path string, problems *Problems) {
	if !node.HasNode(f.name) {
		return
	}

	value := node.GetNode(f.name)

	if f.params.Enum != nil && len(f.params.Enum) != 0 {
		include := false

		for _, param := range f.params.Enum {
			paramDeduce := DeduceValue(param)
			if reflect.DeepEqual(value.Value(), paramDeduce) {
				include = true

				break
			}
		}

		if !include {
			problemStr := localization.GetLocalizerOrPanic().GetString("FieldValueConstraintEnum")
			problems.Append(fmt.Sprintf(problemStr, path, f.name, value.value, f.params.Enum))
		}
		// if enum constraint are provided, others are ignored
		return
	}

	if f.params.IsEqual {
		// if not a number, should be caught by type constraint
		if value.Is(Numeric) && value.ToFloat() != f.params.Equal {
			problemStr := localization.GetLocalizerOrPanic().GetString("FieldValueConstraintEqual")
			problems.Append(fmt.Sprintf(problemStr, path, f.name, value.value, f.params.Equal))
		}
	}

	if f.params.IsFloor {
		// if not a number, should be caught by type constraint
		if value.Is(Numeric) && value.ToFloat() <= f.params.Floor {
			problemStr := localization.GetLocalizerOrPanic().GetString("FieldValueConstraintFloor")
			problems.Append(fmt.Sprintf(problemStr, path, f.name, value.value, f.params.Floor))
		}
	}

	if f.params.IsMin {
		// if not a number, should be caught by type constraint
		if value.Is(Numeric) && value.ToFloat() < f.params.Min {
			problemStr := localization.GetLocalizerOrPanic().GetString("FieldValueConstraintMin")
			problems.Append(fmt.Sprintf(problemStr, path, f.name, value.value, f.params.Min))
		}
	}

	if f.params.IsCeiling {
		// if not a number, should be caught by type constraint
		if value.Is(Numeric) && value.ToFloat() >= f.params.Ceiling {
			problemStr := localization.GetLocalizerOrPanic().GetString("FieldValueConstraintCeiling")
			problems.Append(fmt.Sprintf(problemStr, path, f.name, value.value, f.params.Ceiling))
		}
	}

	if f.params.IsMax {
		// if not a number, should be caught by type constraint
		if value.Is(Numeric) && value.ToFloat() > f.params.Max {
			problemStr := localization.GetLocalizerOrPanic().GetString("FieldValueConstraintMax")
			problems.Append(fmt.Sprintf(problemStr, path, f.name, value.value, f.params.Max))
		}
	}
}
