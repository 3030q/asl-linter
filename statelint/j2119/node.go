package j2119

import (
	"math"
	"time"
)

const (
	Object        = ValueType("object")
	Array         = ValueType("array")
	String        = ValueType("string")
	Integer       = ValueType("integer")
	Float         = ValueType("float")
	Bool          = ValueType("boolean")
	Numeric       = ValueType("numeric")
	JSONPath      = ValueType("JSONPath")
	ReferencePath = ValueType("referencePath")
	Timestamp     = ValueType("timestamp")
	URI           = ValueType("URI")
)

func GetAllTypesString() []string {
	return []string{
		string(Object),
		string(Array),
		string(String),
		string(Integer),
		string(Float),
		string(Bool),
		string(Numeric),
		string(JSONPath),
		string(ReferencePath),
		string(Timestamp),
		string(URI),
	}
}

type ValueType string

type Node struct {
	value      interface{}
	valueTypes map[ValueType]struct{}
}

func NewNode(value interface{}) *Node {
	n := &Node{value: value}
	n.valueTypes = n.getTypes()

	return n
}

func (n *Node) Types() []ValueType {
	result := make([]ValueType, 0, len(n.valueTypes))
	for t := range n.valueTypes {
		result = append(result, t)
	}

	return result
}

func (n *Node) Is(valueType ValueType) bool {
	_, ok := n.valueTypes[valueType]

	return ok
}

func (n *Node) IsNull() bool {
	return n.value == nil
}

func (n *Node) ToTimestamp() time.Time {
	if !n.Is(Timestamp) {
		panic("value type isn't timestamp")
	}

	t, _ := time.Parse(time.RFC3339, n.value.(string))

	return t
}

func (n *Node) ToString() string {
	if !n.Is(String) {
		panic("value type isn't string")
	}

	return n.value.(string)
}

func (n *Node) ToFloat() float64 {
	if !n.Is(Float) {
		panic("value type isn't float")
	}

	if intVal, ok := n.value.(int); ok {
		return float64(intVal)
	}

	return n.value.(float64)
}

func (n *Node) ToObject() map[string]Node {
	if !n.Is(Object) {
		panic("value type isn't object")
	}

	obj, _ := n.value.(map[string]interface{})

	resultObject := make(map[string]Node)
	for key, value := range obj {
		resultObject[key] = *NewNode(value)
	}

	return resultObject
}

func (n *Node) ToInt() int {
	if !n.Is(Integer) {
		panic("value type isn't integer")
	}

	if i, ok := n.value.(int); ok {
		return i
	}

	// Check is not necessary, because it's already checked when creating node
	return int(n.value.(float64))
}

func (n *Node) ToBool() bool {
	if !n.Is(Bool) {
		panic("value type isn't bool")
	}

	return n.value.(bool)
}

func (n *Node) GetNode(name string) *Node {
	if !n.Is(Object) {
		panic("current node is not an object")
	}

	// Check is not necessary, because it's already checked when creating node
	obj, _ := n.value.(map[string]interface{})
	if value, ok := obj[name]; ok {
		return NewNode(value)
	}

	panic("current node not containing node with name " + name)
}

func (n *Node) Keys() []string {
	if !n.Is(Object) {
		panic("node is not an object. Cant get keys")
	}

	obj := n.ToObject()
	keys := make([]string, 0, len(obj))

	for key := range obj {
		keys = append(keys, key)
	}

	return keys
}

func (n *Node) HasNode(name string) bool {
	if !n.Is(Object) {
		return false
	}

	// Check is not necessary, because it's already checked when creating node
	obj, _ := n.value.(map[string]interface{})

	_, ok := obj[name]

	return ok
}

func (n *Node) ValueToArray() []Node {
	if !n.Is(Array) {
		panic("current node is not an array")
	}

	// Check is not necessary, because it's already checked when creating node
	arr, _ := n.value.([]interface{})
	result := make([]Node, 0, len(arr))

	for _, item := range arr {
		result = append(result, *NewNode(item))
	}

	return result
}

// Value returns underlying value, converted to its type
func (n *Node) Value() interface{} {
	switch {
	case n.Is(Array):
		return n.value.([]interface{})
	case n.Is(Object):
		return n.ToObject()
	case n.Is(String):
		return n.ToString()
	case n.Is(Numeric):
		if n.Is(Integer) {
			return n.ToInt()
		}

		return n.ToFloat()
	case n.Is(Bool):
		return n.ToBool()
	default:
		return n.value
	}
}

func (n *Node) getTypes() map[ValueType]struct{} {
	result := make(map[ValueType]struct{})

	switch {
	case n.isArray():
		result[Array] = struct{}{}

		return result
	case n.isObject():
		result[Object] = struct{}{}

		return result
	case n.isString():
		result[String] = struct{}{}

		s, _ := n.value.(string)
		if IsPath(s) {
			result[JSONPath] = struct{}{}
		}

		if IsReferencePath(s) {
			result[ReferencePath] = struct{}{}
		}

		if n.isTimestamp() {
			result[Timestamp] = struct{}{}
		}

		if n.isURI() {
			result[URI] = struct{}{}
		}

		return result
	case n.isNumber():
		result[Numeric] = struct{}{}
		result[Float] = struct{}{}

		if n.isInt() {
			result[Integer] = struct{}{}
		}

		return result
	case n.isBool():
		result[Bool] = struct{}{}

		return result
	default:
		return result
	}
}

func (n *Node) isArray() bool {
	_, ok := n.value.([]interface{})

	return ok
}

func (n *Node) isBool() bool {
	_, ok := n.value.(bool)

	return ok
}

func (n *Node) isInt() bool {
	_, ok := n.value.(int)
	if ok {
		return true
	}

	f, _ := n.value.(float64)

	return math.Mod(f, 1.0) == 0
}

func (n *Node) isNumber() bool {
	_, floatOk := n.value.(float64)
	_, intOk := n.value.(int)

	return floatOk || intOk
}

func (n *Node) isTimestamp() bool {
	value, _ := n.value.(string)

	_, err := time.Parse(time.RFC3339, value)

	return err == nil
}

func (n *Node) isObject() bool {
	_, ok := n.value.(map[string]interface{})

	return ok
}

func (n *Node) isString() bool {
	_, ok := n.value.(string)

	return ok
}

func (n *Node) isURI() bool {
	s, _ := n.value.(string)

	return uriMatcher.MatchString(s)
}
