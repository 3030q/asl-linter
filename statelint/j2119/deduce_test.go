package j2119

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeduce_SpotTypesCorrectly(t *testing.T) {
	t.Parallel()

	assert.EqualValues(t, "foo", DeduceValue(`"foo"`))
	assert.EqualValues(t, "foo", DeduceValue(`foo`))
	assert.EqualValues(t, true, DeduceValue(`true`))
	assert.EqualValues(t, false, DeduceValue(`false`))
	assert.EqualValues(t, nil, DeduceValue(`null`))
	assert.EqualValues(t, 234, DeduceValue(`234`))
	assert.EqualValues(t, 25.411, DeduceValue(`25.411`))
}
