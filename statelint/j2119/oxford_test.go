package j2119

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOxford_PatternWorking(t *testing.T) {
	t.Parallel()

	re := regexp.MustCompile("^" + basic + "$")
	assert.True(t, re.MatchString("X"))
	assert.True(t, re.MatchString("X or X"))
	assert.True(t, re.MatchString("X, X, or X"))
	assert.True(t, re.MatchString("X, X, X, or X"))
}

func TestOxford_SimpleMatch(t *testing.T) {
	t.Parallel()

	cases := []string{"a",
		"a or aa",
		"a, aa, or aaa",
		"a, aa, aaa, or aaaa"}

	re := regexp.MustCompile("^" + GetOxfordRe("a+", OxfordOptions{}) + "$")

	for _, c := range cases {
		assert.True(t, re.MatchString(c))
	}
}

func TestOxford_UseCaptureArticleConnector(t *testing.T) {
	t.Parallel()

	cases := []string{
		`an "asdg"`,
		`a "foij2pe" and an "aiepw"`,
		`an "alkvm 2", an "ap89wf", and a " lfdj a fddalfkj"`,
		`an "aj89peww", a "", an "aslk9 ", and an "x"`,
	}

	ox := GetOxfordRe(`"([^"]*)"`, OxfordOptions{
		UseArticle:     true,
		CaptureName:    "capture_me",
		HasCaptureName: true,
		Connector:      "and",
		HasConnector:   true,
	})
	re := regexp.MustCompile("^" + ox + "$")

	for _, c := range cases {
		assert.True(t, re.MatchString(c))
	}
}

func TestOxford_BreakUpRoleList(t *testing.T) {
	t.Parallel()

	list := []string{
		"an R2",
		"an R2 or an R3",
		"an R2, an R3, or an R4",
	}

	wantedPieces := [][]string{
		{"R2"},
		{"R2", "R3"},
		{"R2", "R3", "R4"},
	}

	matcher := NewMatcher("R1")
	for i := 2; i < 5; i++ {
		matcher.AddRole(fmt.Sprintf("R%d", i))
	}

	for i, l := range list {
		resultedList := BreakRoleList(*matcher, l)
		isEqual := reflect.DeepEqual(resultedList, wantedPieces[i])
		assert.True(t, isEqual)
	}
}

func TestOxford_BreakUpStringList(t *testing.T) {
	t.Parallel()

	cases := []string{
		`"R2"`,
		`"R2" or "R3"`,
		`"R2", "R3", or "R4"`,
	}

	wantedPieces := [][]string{
		{"R2"},
		{"R2", "R3"},
		{"R2", "R3", "R4"},
	}

	for i, c := range cases {
		resultedList := BreakStringList(c)
		isEqual := reflect.DeepEqual(resultedList, wantedPieces[i])
		assert.True(t, isEqual)
	}
}
