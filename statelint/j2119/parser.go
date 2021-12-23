package j2119

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
)

var (
	rootRegex   = regexp.MustCompile(`This\s+document\s+specifies\s+a\s+JSON\s+object\s+called\s+an?\s+"([^"]+)"\.`)
	eachOfRegex = regexp.MustCompile(`^Each of a`)
)

type Parser interface {
	FindMoreRoles(node Node, roles []string) []string
	FindGrandchildRoles(roles []string, name string) []string
	FindChildRoles(roles []string, name string) []string
	GetConstraints(role string) []Constrainter
	IsFieldAllowed(roles []string, child string) bool
	IsAllowsAny(roles []string) bool
}

type J2119Parser struct {
	haveRoot      bool
	root          string
	constraints   *RoleConstraints
	finder        *RoleFinder
	allowedFields *AllowedFields
	matcher       *Matcher
	assigner      *Assigner
}

func NewJ2119Parser(j2119File io.Reader) (*J2119Parser, error) {
	p := &J2119Parser{
		constraints:   NewRoleConstraints(),
		finder:        NewRoleFinder(),
		allowedFields: NewAllowedFields(),
	}
	scanner := bufio.NewScanner(j2119File)

	for scanner.Scan() {
		line := scanner.Text()

		if rootRegex.MatchString(line) {
			if p.haveRoot {
				panic("Only one root declaration")
			}

			p.root = rootRegex.FindStringSubmatch(line)[1]
			p.matcher = NewMatcher(p.root)
			p.assigner = NewAssigner(p.constraints, p.finder, p.matcher, p.allowedFields)
			p.haveRoot = true

			continue
		}

		if !p.haveRoot {
			panic("Root declaration must go first")
		}

		err := p.procLine(line)

		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *J2119Parser) procLine(line string) error {
	switch {
	case p.matcher.IsConstraintLine(line):
		p.assigner.AssignConstraints(p.matcher.BuildConstraint(line))
	case p.matcher.IsOnlyOneMatchLine(line):
		p.assigner.AssignOnlyOneOf(p.matcher.BuildOnlyOne(line))
	case eachOfRegex.MatchString(line):
		eachesLine := p.matcher.BuildEachOfLine(line)
		eaches := BreakRoleList(*p.matcher, eachesLine["each_of"])

		for _, each := range eaches {
			err := p.procLine(fmt.Sprintf("A %s %s", each, eachesLine["trailer"]))
			if err != nil {
				return err
			}
		}
	case p.matcher.IsRoleDefLine(line):
		p.assigner.AssignRoles(p.matcher.BuildRoleDef(line))
	default:
		return errors.New(fmt.Sprintf("Unrecognized Line: %s", line))
	}

	return nil
}

func (p *J2119Parser) FindMoreRoles(node Node, roles []string) []string {
	return p.finder.FindMoreRoles(node, roles)
}

func (p *J2119Parser) FindGrandchildRoles(roles []string, name string) []string {
	return p.finder.FindGrandchildRoles(roles, name)
}

func (p *J2119Parser) FindChildRoles(roles []string, name string) []string {
	return p.finder.FindChildRoles(roles, name)
}

func (p *J2119Parser) GetConstraints(role string) []Constrainter {
	return p.constraints.Get(role)
}

func (p *J2119Parser) IsFieldAllowed(roles []string, child string) bool {
	return p.allowedFields.IsAllowed(roles, child)
}

func (p *J2119Parser) IsAllowsAny(roles []string) bool {
	return p.allowedFields.IsAny(roles)
}
