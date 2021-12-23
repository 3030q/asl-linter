package j2119

import (
	"fmt"
	"statelint/localization"
)

type NodeValidator struct {
	parser Parser
}

type NodeValidatorParams struct {
	node  Node
	path  string
	roles []string
}

func NewNodeValidatorParams(node Node, path string, roles []string) *NodeValidatorParams {
	return &NodeValidatorParams{node: node, path: path, roles: roles}
}

func NewNodeValidator(parser Parser) *NodeValidator {
	return &NodeValidator{parser: parser}
}

func (n *NodeValidator) Validate(initialNode Node, initialPath string, initialRoles []string, problems *Problems) {
	nodes := []*NodeValidatorParams{NewNodeValidatorParams(initialNode, initialPath, initialRoles)}

	for len(nodes) != 0 {
		currentNode := nodes[0]
		nodes = nodes[1:]

		currentNode.roles = n.parser.FindMoreRoles(currentNode.node, currentNode.roles)

		for _, role := range currentNode.roles {
			for _, constraint := range n.parser.GetConstraints(role) {
				if constraint.Applies(currentNode.node, currentNode.roles) {
					constraint.Check(currentNode.node, currentNode.path, problems)
				}
			}
		}

		if !currentNode.node.Is(Object) {
			continue
		}

		for name, val := range currentNode.node.ToObject() {
			if !n.parser.IsFieldAllowed(currentNode.roles, name) {
				problemStr := localization.GetLocalizerOrPanic().GetString("NodeValidatorIsFieldAllowed")
				problems.Append(fmt.Sprintf(problemStr, name, currentNode.path))
			}

			// only recurse into children if they have roles
			childRoles := n.parser.FindChildRoles(currentNode.roles, name)
			if len(childRoles) != 0 {
				nodes = append(nodes, NewNodeValidatorParams(
					val,
					fmt.Sprintf("%s.%s", currentNode.path, name),
					childRoles,
				))
			}

			grandchildRoles := n.parser.FindGrandchildRoles(currentNode.roles, name)
			if len(grandchildRoles) != 0 && !n.parser.IsAllowsAny(grandchildRoles) {
				// recurse into grandkids
				switch {
				case val.Is(Object):
					for childName, childVal := range val.ToObject() {
						nodes = append(nodes, NewNodeValidatorParams(
							childVal,
							fmt.Sprintf("%s.%s.%s", currentNode.path, name, childName),
							grandchildRoles,
						))
					}
				case val.Is(Array):
					for i, member := range val.ValueToArray() {
						nodes = append(nodes, NewNodeValidatorParams(
							member,
							fmt.Sprintf("%s.%s[%d]", currentNode.path, name, i),
							grandchildRoles,
						))
					}
				}
			}
		}
	}
}
