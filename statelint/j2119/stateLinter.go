package j2119

const DefaultStateMachinePath = "./data/StateMachine.j2119"

type StateLinter struct {
	validator *Validator
}

func NewStateLinter() (*StateLinter, error) {
	return NewStateLinterFromFile(DefaultStateMachinePath)
}

func NewStateLinterFromFile(path string) (*StateLinter, error) {
	v, err := NewValidator(path)
	if err != nil {
		return nil, err
	}

	return &StateLinter{
		validator: v,
	}, nil
}

func (s *StateLinter) ValidateJSONStruct(jsonObject interface{}) *Problems {
	node := *NewNode(jsonObject)
	problems := s.validator.ValidateJSONStruct(jsonObject)

	// additional check
	stateNode := NewStateNode()
	stateNode.Check(node, s.validator.parser.root, problems)

	return problems
}
