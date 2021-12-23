package j2119

type AllowedFields struct {
	allowed map[string]map[string]struct{}
	any     map[string]struct{}
}

func NewAllowedFields() *AllowedFields {
	return &AllowedFields{
		allowed: make(map[string]map[string]struct{}),
		any:     make(map[string]struct{}),
	}
}

func (a *AllowedFields) SetAllowed(role string, child string) {
	if _, ok := a.allowed[role]; !ok {
		a.allowed[role] = make(map[string]struct{})
	}

	a.allowed[role][child] = struct{}{}
}

func (a *AllowedFields) SetAny(role string) {
	a.any[role] = struct{}{}
}

func (a *AllowedFields) IsAllowed(roles []string, child string) bool {
	if a.IsAny(roles) || len(roles) != 0 {
		for _, role := range roles {
			if _, hasRole := a.allowed[role]; hasRole {
				if _, hasChild := a.allowed[role][child]; hasChild {
					return true
				}
			}
		}
	}

	return false
}

func (a *AllowedFields) IsAny(roles []string) bool {
	for _, role := range roles {
		if _, has := a.any[role]; has {
			return true
		}
	}

	return false
}
