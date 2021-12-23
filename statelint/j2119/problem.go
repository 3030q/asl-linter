package j2119

type Problems struct {
	problems []string
}

func NewProblems() *Problems {
	return &Problems{
		problems: []string{},
	}
}

func (p *Problems) Append(value string) {
	p.problems = append(p.problems, value)
}

func (p *Problems) Len() int {
	if p.problems == nil {
		panic("problems array not initialized")
	}

	return len(p.problems)
}

func (p *Problems) GetProblems() []string {
	return p.problems
}
