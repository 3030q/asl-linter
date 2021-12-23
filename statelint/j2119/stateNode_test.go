package j2119

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetNodeFromTestFile(t *testing.T, filename string) Node {
	t.Helper()

	p := "../testdata/" + filename

	openedFile, err := os.Open(p)
	if err != nil {
		t.Fatalf("can not open file %s", p)
	}
	defer openedFile.Close()

	all, err := io.ReadAll(openedFile)
	if err != nil {
		t.Fatalf("can not read file %s", p)
	}

	return NewNodeCreateHelper(t, string(all))
}

func TestStateNode_FindMissingStartAtTargets(t *testing.T) {
	t.Parallel()

	json := `{
			   "StartAt": "x",
			   "States": {
				 "y": {
				   "Type": "Succeed"
				 }
			   }
			 }`

	node := NewNodeCreateHelper(t, json)
	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 2, problems.Len())
}

func TestStateNode_CatchNestedProblems(t *testing.T) {
	t.Parallel()

	json := `{
			  "StartAt": "x",
			  "States": {
				"x": {
				  "StartAt": "z",
				  "States": {
					"w": 1
				  }
				}
			  }
			}`

	node := NewNodeCreateHelper(t, json)
	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 4, problems.Len())
}

func TestStateNode_FindStatesALLNotInEndPosition(t *testing.T) {
	t.Parallel()

	json := `{
			  "Retry": [
				{
				  "ErrorEquals": [
					"States.ALL",
					"other"
				  ]
				},
				{
				  "ErrorEquals": [
					"YET ANOTHER"
				  ]
				}
			  ]
			}`

	node := NewNodeCreateHelper(t, json)
	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestStateNode_FindStatesALLNotByItself(t *testing.T) {
	t.Parallel()

	json := `{
			  "Retry": [
				{
				  "ErrorEquals": [
					"YET ANOTHER"
				  ]
				},
				{
				  "ErrorEquals": [
					"States.ALL",
					"other"
				  ]
				}
			  ]
			}`

	node := NewNodeCreateHelper(t, json)
	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestStateNode_UseDefaultFieldCorrectly(t *testing.T) {
	t.Parallel()

	json := `{
			  "StartAt": "A",
			  "States": {
				"A": {
				  "Type": "Choice",
				  "Choices": [
					{
					  "Variable": "$.a",
					  "Next": "B"
					}
				  ],
				  "Default": "C"
				},
				"B": {
				  "Type": "Succeed"
				},
				"C": {
				  "Type": "Succeed"
				}
			  }
			}`

	node := NewNodeCreateHelper(t, json)
	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 0, problems.Len())
}

func TestStateNode_FindNextFieldsWithTargetsThatDontMatchStateNames(t *testing.T) {
	t.Parallel()

	json := `{
			  "StartAt": "A",
			  "States": {
				"A": {
				  "Type": "Pass",
				  "Next": "B"
				}
			  }
			}`

	node := NewNodeCreateHelper(t, json)
	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 2, problems.Len())
}

func TestStateNode_ShouldFindUnPointedToStates(t *testing.T) {
	t.Parallel()

	json := `{
			  "StartAt": "A",
			  "States": {
				"A": {
				  "Type": "Succeed"
				},
				"X": {
				  "Type": "Succeed"
				}
			  }
			}`

	node := NewNodeCreateHelper(t, json)
	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestStateNode_FindMissingTerminalState(t *testing.T) {
	t.Parallel()

	json := `{
			  "StartAt": "A",
			  "States": {
				"A": {
				  "Type": "Pass",
				  "Next": "B"
				},
				"B": {
				  "Type": "C",
				  "Next": "A"
				}
			  }
			}`

	node := NewNodeCreateHelper(t, json)
	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestStateNode_HandleComplexMissingTerminal(t *testing.T) {
	t.Parallel()

	node := GetNodeFromTestFile(t, "noTerminal.json")

	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 1, problems.Len())
}

func TestStateNode_CatchLinkageFromOneParallelBranchToAnother(t *testing.T) {
	t.Parallel()

	node := GetNodeFromTestFile(t, "linkedParallel.json")

	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 4, problems.Len())
}

func TestStateNode_CatchDuplicatedStateNamesEvenInParallels(t *testing.T) {
	t.Parallel()

	node := GetNodeFromTestFile(t, "hasDupes.json")

	problems := NewProblems()
	checker := NewStateNode()
	checker.Check(node, "a.b", problems)

	assert.Equal(t, 1, problems.Len())
}
