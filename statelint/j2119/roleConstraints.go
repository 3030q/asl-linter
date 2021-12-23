package j2119

type RoleConstraints struct {
	constraints map[string][]Constrainter
}

func NewRoleConstraints() *RoleConstraints {
	return &RoleConstraints{
		constraints: make(map[string][]Constrainter),
	}
}

func (r *RoleConstraints) Add(role string, constraint Constrainter) {
	r.constraints[role] = append(r.constraints[role], constraint)
}

func (r *RoleConstraints) Get(role string) []Constrainter {
	if _, ok := r.constraints[role]; !ok {
		return []Constrainter{}
	}

	return r.constraints[role]
}
