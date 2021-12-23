package j2119

import (
	"fmt"
	"regexp"
	"statelint/localization"
	"strings"
)

var intrinsicInvocationRegex = regexp.MustCompile(`^States\.(JsonToString|Format|StringToJson|Array)\(.+\)$`)

type StateNode struct {
	currentStatesNode     []Node
	currentStatesIncoming [][]string

	allStateNames              map[string]string
	payloadBuilderFields       []string
	contextObjectAccessField   []map[string]interface{}
	choiceStateNestedOperators []string
}

func NewStateNode() *StateNode {
	return &StateNode{
		currentStatesNode:     make([]Node, 0),
		currentStatesIncoming: make([][]string, 0),
		allStateNames:         make(map[string]string),
		payloadBuilderFields:  []string{"Parameters", "ResultSelector"},
		contextObjectAccessField: []map[string]interface{}{
			{"field": "InputPath", "nullable": true},
			{"field": "OutputPath", "nullable": true},
			{"field": "ItemsPath", "nullable": false},
		},
		choiceStateNestedOperators: []string{"And", "Or", "Not"},
	}
}

func (s *StateNode) Check(node Node, path string, problems *Problems) {
	if !node.Is(Object) {
		return
	}

	isMachineTop := node.HasNode("States") && node.GetNode("States").Is(Object)

	if isMachineTop {
		s.currentStatesNode = append(s.currentStatesNode, *node.GetNode("States"))

		if node.HasNode("StartAt") && node.GetNode("StartAt").Is(String) {
			startAt := node.GetNode("StartAt").ToString()
			s.currentStatesIncoming = append(s.currentStatesIncoming, []string{startAt})

			if !node.GetNode("States").HasNode(startAt) {
				problemStr := localization.GetLocalizerOrPanic().
					GetString("StateNodeDoesntHaveStartAtNode")
				problems.Append(fmt.Sprintf(problemStr, startAt, path))
			}
		} else {
			s.currentStatesIncoming = append(s.currentStatesIncoming, []string{})
		}

		states := node.GetNode("States")
		for _, name := range states.Keys() {
			child := *states.GetNode(name)

			if child.Is(Object) {
				childPath := path + "." + name
				s.ProbeContextObjectAccess(child, childPath, problems)

				for _, fieldName := range s.payloadBuilderFields {
					if child.HasNode(fieldName) {
						s.ProbePayloadBuilder(
							*child.GetNode(fieldName),
							childPath,
							problems,
							fieldName,
						)
					}
				}

				if child.HasNode("Type") &&
					child.GetNode("Type").Is(String) &&
					child.GetNode("Type").ToString() == "Choice" &&
					child.HasNode("Choices") {
					s.ProbeChoiceState(*child.GetNode("Choices"), childPath+".Choices", problems)
				}
			}

			if _, ok := s.allStateNames[name]; ok {
				problemStr := localization.GetLocalizerOrPanic().
					GetString("StateNodeDoubleDefinedState")
				problems.Append(fmt.Sprintf(problemStr, name, path, s.allStateNames[name]))
			} else {
				s.allStateNames[name] = fmt.Sprintf("%s.States", path)
			}
		}
	}

	s.CheckForTerminal(node, path, problems)

	s.CheckNext(node, path, problems)

	if node.HasNode("Retry") {
		s.checkStatesAll(*node.GetNode("Retry"), path+".Retry", problems)
	}

	if node.HasNode("Catch") {
		s.checkStatesAll(*node.GetNode("Catch"), path+".Catch", problems)
	}

	for _, name := range node.Keys() {
		val := *node.GetNode(name)

		if val.Is(Array) {
			for i, element := range val.ValueToArray() {
				s.Check(element, fmt.Sprintf("%s.%s[%d]", path, name, i), problems)
			}
		} else {
			s.Check(val, fmt.Sprintf("%s.%s", path, name), problems)
		}
	}

	if isMachineTop {
		states := s.currentStatesNode[len(s.currentStatesNode)-1]
		s.currentStatesNode = s.currentStatesNode[:len(s.currentStatesNode)-1]
		incoming := s.currentStatesIncoming[len(s.currentStatesIncoming)-1]
		s.currentStatesIncoming = s.currentStatesIncoming[:len(s.currentStatesIncoming)-1]

		stateKeys := states.Keys()

		var missing []string

		for _, key := range stateKeys {
			has := false

			for _, incomingKey := range incoming {
				if key == incomingKey {
					has = true

					break
				}
			}

			if !has {
				missing = append(missing, key)
			}
		}

		for _, state := range missing {
			problemStr := localization.GetLocalizerOrPanic().GetString("StateNodeMissingTransition")
			problems.Append(fmt.Sprintf(problemStr, path, state))
		}
	}
}

func (s *StateNode) CheckNext(node Node, path string, problems *Problems) {
	s.AddNext(node, path, "Next", problems)
	s.AddNext(node, path, "Default", problems)
}

func (s *StateNode) AddNext(node Node, path string, field string, problems *Problems) {
	if !node.HasNode(field) || !node.GetNode(field).Is(String) {
		return
	}

	transitionTo := node.GetNode(field).ToString()

	if len(s.currentStatesNode) != 0 {
		if s.currentStatesNode[len(s.currentStatesNode)-1].HasNode(transitionTo) {
			lastIndex := len(s.currentStatesIncoming) - 1
			s.currentStatesIncoming[lastIndex] =
				append(s.currentStatesIncoming[lastIndex], transitionTo)
		} else {
			problemStr := localization.GetLocalizerOrPanic().GetString("StateNodeNoStateFoundReferenced")
			problems.Append(fmt.Sprintf(problemStr, transitionTo, path, field))
		}
	}
}

func (s *StateNode) ProbeContextObjectAccess(node Node, path string, problems *Problems) {
	for _, field := range s.contextObjectAccessField {
		fieldName, _ := field["field"].(string)
		nullable, _ := field["nullable"].(bool)

		if !node.HasNode(fieldName) {
			continue
		}

		if !nullable && node.GetNode(fieldName).IsNull() {
			problemStr := localization.GetLocalizerOrPanic().GetString("StateNodeFieldShouldBeNonNull")
			problems.Append(fmt.Sprintf(problemStr, fieldName, path))

			return
		}

		if !node.GetNode(fieldName).IsNull() && !s.IsValidParametersPath(*node.GetNode(fieldName)) {
			problemStr := localization.GetLocalizerOrPanic().GetString("StateNodeFieldIsNotJSONPath")
			problems.Append(fmt.Sprintf(problemStr, fieldName, path))
		}
	}
}

func (s *StateNode) ProbeChoiceState(node Node, path string, problems *Problems) {
	switch {
	case node.Is(Object):
		if node.HasNode("Variable") && !s.IsValidParametersPath(*node.GetNode("Variable")) {
			problemStr := localization.GetLocalizerOrPanic().GetString("StateNodeProbChoiceState")
			problems.Append(fmt.Sprintf(problemStr, path))
		}

		for _, operator := range s.choiceStateNestedOperators {
			if node.HasNode(operator) {
				s.ProbeChoiceState(
					*node.GetNode(operator),
					fmt.Sprintf("%s.%s", path, operator),
					problems,
				)
			}
		}
	case node.Is(Array):
		arrayNodes := node.ValueToArray()
		for i, arrayNode := range arrayNodes {
			s.ProbeChoiceState(arrayNode, fmt.Sprintf("%s[%d]", path, i), problems)
		}
	}
}

func (s *StateNode) ProbePayloadBuilder(node Node, path string, problems *Problems, fieldName string) {
	switch {
	case node.Is(Object):
		for _, key := range node.Keys() {
			value := *node.GetNode(key)
			if strings.HasSuffix(key, ".$") {
				if !s.IsIntrinsicInvocation(value) && !s.IsValidParametersPath(value) {
					problemStr := localization.GetLocalizerOrPanic().
						GetString("StateNodeProbePayloadBuilder")
					problems.Append(fmt.Sprintf(problemStr, key, fieldName, path))
				}

				continue
			}

			s.ProbePayloadBuilder(
				value,
				fmt.Sprintf("%s.%s", path, key),
				problems,
				fieldName,
			)
		}
	case node.Is(Array):
		arrayNodes := node.ValueToArray()
		for i, arrayNode := range arrayNodes {
			s.ProbePayloadBuilder(arrayNode, fmt.Sprintf("%s[%d]", path, i), problems, fieldName)
		}
	}
}

func (s *StateNode) IsIntrinsicInvocation(value Node) bool {
	return value.Is(String) && intrinsicInvocationRegex.MatchString(value.ToString())
}

func (s *StateNode) IsValidParametersPath(value Node) bool {
	if !value.Is(String) {
		return false
	}

	stringValue := value.ToString()

	if strings.HasPrefix(stringValue, "$$") {
		stringValue = stringValue[1:]

		return IsPath(stringValue)
	}

	return IsReferencePath(stringValue)
}

var terminalTypes = []interface{}{"Succeed", "Fail"}

func (s *StateNode) CheckForTerminal(node Node, path string, problems *Problems) {
	var terminalFound bool

	if !node.HasNode("States") || !node.GetNode("States").Is(Object) {
		return
	}

	states := node.GetNode("States")

	for _, key := range states.Keys() {
		stateNode := states.GetNode(key)
		if stateNode.HasNode("Type") && stateNode.GetNode("Type").Is(String) {
			typeNode := stateNode.GetNode("Type").ToString()

			for _, t := range terminalTypes {
				if t == typeNode {
					terminalFound = true

					break
				}
			}
		}

		if stateNode.HasNode("End") {
			terminalFound = true
		}

		if terminalFound {
			break
		}
	}

	if !terminalFound {
		problemStr := localization.GetLocalizerOrPanic().GetString("StateNodeCheckForTerminal")
		problems.Append(fmt.Sprintf(problemStr, path))
	}
}

func (s *StateNode) checkStatesAll(node Node, path string, problems *Problems) {
	if !node.Is(Array) {
		return
	}

	nodes := node.ValueToArray()
	for i, element := range nodes {
		if element.Is(Object) &&
			element.HasNode("ErrorEquals") &&
			element.GetNode("ErrorEquals").Is(Array) {
			ee := element.GetNode("ErrorEquals").ValueToArray()
			has := false

			for _, eeElement := range ee {
				if eeElement.Is(String) && eeElement.ToString() == "States.ALL" {
					has = true

					break
				}
			}

			if has && (i != len(nodes)-1 || len(ee) != 1) {
				problemStr := localization.GetLocalizerOrPanic().GetString("StateNodeCheckStatesAll")
				problems.Append(fmt.Sprintf(problemStr, path, i))
			}
		}
	}
}
