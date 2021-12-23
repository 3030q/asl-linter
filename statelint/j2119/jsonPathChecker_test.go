package j2119

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONPathChecker_AllowDefaultPaths(t *testing.T) {
	t.Parallel()

	assert.True(t, IsPath(`$`))
	assert.True(t, IsReferencePath(`$`))
}

func TestJSONPathChecker_DoSimplePaths(t *testing.T) {
	t.Parallel()

	assert.True(t, IsPath(`$.foo.bar`))
	assert.True(t, IsPath(`$..x`))
	assert.True(t, IsPath(`$.foo.bar.baz.biff..blecch`))
	assert.True(t, IsPath(`$.café_au_lait`))
	assert.True(t, IsPath(`$['foo']`))
	assert.True(t, IsPath(`$[3]`))
}

func TestJSONPathChecker_RejectObviousBotches(t *testing.T) {
	t.Parallel()

	assert.False(t, IsPath(`x`))
	assert.False(t, IsPath(`.x`))
	assert.False(t, IsPath(`x.y.z`))
	assert.False(t, IsPath(`$.~.bar`))
	assert.False(t, IsReferencePath(`x`))
	assert.False(t, IsReferencePath(`.x`))
	assert.False(t, IsReferencePath(`x.y.z`))
	assert.False(t, IsReferencePath(`$.~.bar`))
}

func TestJSONPathChecker_AcceptPathsWithBracketNotation(t *testing.T) {
	t.Parallel()

	assert.True(t, IsPath(`$['foo']['bar']`))
	assert.True(t, IsPath(`$['foo']['bar']['baz']['biff']..blecch`))
	assert.True(t, IsPath(`$['café_au_lait']`))
}

func TestJSONPathChecker_AcceptJaywayJSONPathExamples(t *testing.T) {
	t.Parallel()

	paths := []string{
		`$.store.book[*].author`,
		`$..author`,
		`$.store.*`,
		`$..book[2]`,
		`$..book[0,1]`,
		`$..book[:2]`,
		`$..book[1:2]`,
		`$..book[-2:]`,
		`$..book[2:]`,
		`$..*`,
	}

	for _, p := range paths {
		assert.True(t, IsPath(p))
	}
}

func TestJSONPathChecker_AllowReferencePaths(t *testing.T) {
	t.Parallel()

	paths := []string{
		`$.foo.bar`,
		`$..x`,
		`$.foo.bar.baz.biff..blecch`,
		`$.café_au_lait`,
		`$['foo']['bar']`,
		`$['foo']['bar']['baz']['biff']..blecch`,
		`$['café_au_lait']`,
		`$..author`,
		`$..book[2]`,
	}
	for _, p := range paths {
		assert.True(t, IsReferencePath(p))
	}
}

func TestJSONPathChecker_DistinguishBetweenNonPathsPathsAndReferencePaths(t *testing.T) {
	t.Parallel()

	paths := []string{
		`$.store.book[*].author`,
		`$..author`,
		`$.store.*`,
		`$..book[2]`,
		`$..book[0,1]`,
		`$..book[:2]`,
		`$..book[1:2]`,
		`$..book[-2:]`,
		`$..book[2:]`,
		`$..*`,
	}
	referencePaths := []string{
		`$..author`,
		`$..book[2]`,
	}

	for _, p := range paths {
		assert.True(t, IsPath(p))

		include := false

		for _, rp := range referencePaths {
			if p == rp {
				include = true

				break
			}
		}

		if include {
			assert.True(t, IsReferencePath(p))
		} else {
			assert.False(t, IsReferencePath(p))
		}
	}
}
