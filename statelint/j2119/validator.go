package j2119

import (
	"fmt"
	"os"
)

type Validator struct {
	parser      *J2119Parser
	currentNode *Node
}

func NewValidator(assertionSourcePath string) (*Validator, error) {
	assertion, err := os.Open(assertionSourcePath)
	if err != nil {
		return nil, fmt.Errorf("can not open assertion file \"%s\": %w", assertionSourcePath, err)
	}

	parser, err := NewJ2119Parser(assertion)
	if err != nil {
		return nil, err
	}

	return &Validator{parser: parser}, nil
}

func (v *Validator) ValidateJSONStruct(json interface{}) *Problems {
	v.currentNode = NewNode(json)
	problems := NewProblems()

	validator := NewNodeValidator(v.parser)
	validator.Validate(*v.currentNode, v.parser.root, []string{v.parser.root}, problems)

	return problems
}
