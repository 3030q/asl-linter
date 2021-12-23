package j2119

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllowedFields_TrueAnswer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		roles []string
		child string
	}{
		{
			roles: []string{"foo"},
			child: "bar",
		},
		{
			roles: []string{"bar", "baz", "foo"},
			child: "bar",
		},
	}

	cut := NewAllowedFields()
	cut.SetAllowed("foo", "bar")

	for _, testCase := range testCases {
		assert.True(t, cut.IsAllowed(testCase.roles, testCase.child),
			fmt.Sprintf("test case %v should be true", testCase))
	}
}

func TestAllowedFields_FalseAnswers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		roles []string
		child string
	}{
		{
			roles: []string{"foo"},
			child: "baz",
		},
		{
			roles: []string{"bar", "baz", "foo"},
			child: "baz",
		},
	}

	cut := NewAllowedFields()
	cut.SetAllowed("foo", "bar")

	for _, testCase := range testCases {
		assert.False(t, cut.IsAllowed(testCase.roles, testCase.child),
			fmt.Sprintf("test case %v should be false", testCase))
	}
}

func TestAllowedFields_EdgeCases(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		roles []string
		child string
	}{
		{
			roles: []string{"boo"},
			child: "baz",
		},
		{
			roles: []string{},
			child: "baz",
		},
	}

	cut := NewAllowedFields()
	cut.SetAllowed("foo", "bar")

	for _, testCase := range testCases {
		assert.False(t, cut.IsAllowed(testCase.roles, testCase.child),
			fmt.Sprintf("test case %v should be fasle", testCase))
	}
}
