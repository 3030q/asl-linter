package j2119

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditional_FailOnAnExecuteRole(t *testing.T) {
	t.Parallel()

	c := NewRoleNotPresentedCondition([]string{"foo", "bar"})

	var j interface{}
	err := json.Unmarshal([]byte(`{ "bar": 1 }`), &j)

	if err != nil {
		t.Fatal("cant unmarshal json")
	}

	node := NewNode(j)

	assert.False(t, c.IsConstraintApplies(*node, []string{"foo"}))
}

func TestConditional_SucceedOnNonExcludeRole(t *testing.T) {
	t.Parallel()

	c := NewRoleNotPresentedCondition([]string{"foo", "bar"})

	var j interface{}
	err := json.Unmarshal([]byte(`{ "bar": 1 }`), &j)

	if err != nil {
		t.Fatal("cant unmarshal json")
	}

	node := NewNode(j)
	assert.True(t, c.IsConstraintApplies(*node, []string{"baz"}))
}
