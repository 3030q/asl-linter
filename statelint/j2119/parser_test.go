package j2119

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_MatchROOT(t *testing.T) {
	t.Parallel()

	assert.True(t, rootRegex.MatchString(`This document specifies a JSON object called a "State Machine".`))
}

func GetParserAndValidator(t *testing.T) (*J2119Parser, *NodeValidator) {
	t.Helper()

	f, err := os.Open("../data/AWL.j2119")
	if err != nil {
		t.Fatalf("cant open j2119 file, %s", err.Error())
	}

	p, err := NewJ2119Parser(f)
	if err != nil {
		t.Fatalf("error when parsing j2119 file, %s", err.Error())
	}

	v := NewNodeValidator(p)

	return p, v
}

func TestParser_Explore(t *testing.T) {
	t.Parallel()

	p, v := GetParserAndValidator(t)

	obj := make(map[string]interface{})
	obj["StartAt"] = "pass"
	runParserTest(t, v, p, obj, 1)

	delete(obj, "StartAt")

	states := make(map[string]interface{})
	obj["States"] = states
	runParserTest(t, v, p, obj, 1)

	obj["StartAt"] = "pass"
	runParserTest(t, v, p, obj, 0)

	obj["Version"] = 3
	runParserTest(t, v, p, obj, 1)

	obj["Version"] = "1.0"
	obj["Comment"] = true
	runParserTest(t, v, p, obj, 1)

	obj["Comment"] = "Hi"
	pass := make(map[string]interface{})
	states["Pass"] = pass

	runParserTest(t, v, p, obj, 2)

	pass["Next"] = "s1"
	pass["Type"] = "Pass"

	runParserTest(t, v, p, obj, 0)

	pass["Type"] = "flibber"

	runParserTest(t, v, p, obj, 1)

	pass["Type"] = "Pass"
	pass["Comment"] = 23.5

	runParserTest(t, v, p, obj, 1)

	pass["Type"] = "Pass"
	pass["Comment"] = ""
	pass["End"] = 11

	runParserTest(t, v, p, obj, 1)

	pass["End"] = true
	runParserTest(t, v, p, obj, -1)

	delete(pass, "Next")
	runParserTest(t, v, p, obj, 0)

	pass["InputPath"] = 1
	pass["ResultPath"] = 2

	runParserTest(t, v, p, obj, 2)

	pass["InputPath"] = "foo"
	pass["ResultPath"] = "bar"

	runParserTest(t, v, p, obj, 2)

	fail := map[string]interface{}{
		"Type":  "Fail",
		"Cause": "a",
		"Error": "b",
	}
	states["Fail"] = fail

	delete(pass, "InputPath")
	delete(pass, "ResultPath")

	runParserTest(t, v, p, obj, 0)

	fail["InputPath"], fail["ResultPath"] = "foo", "foo"

	runParserTest(t, v, p, obj, 3)

	delete(fail, "InputPath")
	delete(fail, "ResultPath")
	runParserTest(t, v, p, obj, 0)

	fail["Cause"] = false

	runParserTest(t, v, p, obj, 1)

	fail["Cause"] = "ouch"

	runParserTest(t, v, p, obj, 0)

	task := map[string]interface{}{
		"Type":     "Task",
		"Resource": "a:b",
		"Next":     "fail",
	}
	states["Task"] = task

	runParserTest(t, v, p, obj, 0)

	task["End"] = true
	delete(task, "Next")
	runParserTest(t, v, p, obj, 0)

	task["Resource"] = "foo:bar"
	task["TimeoutSeconds"] = "x"
	task["HeartbeatSeconds"] = 3.9

	runParserTest(t, v, p, obj, -1)

	task["TimeoutSeconds"] = -2
	task["HeartbeatSeconds"] = 0

	runParserTest(t, v, p, obj, -1)

	task["TimeoutSeconds"] = 33
	task["HeartbeatSeconds"] = 44

	runParserTest(t, v, p, obj, 0)

	task["Retry"] = 1

	runParserTest(t, v, p, obj, 1)

	task["Retry"] = []interface{}{1}

	runParserTest(t, v, p, obj, 1)

	task["Retry"] = []interface{}{
		map[string]interface{}{"MaxAttempts": 0},
		map[string]interface{}{"BackoffRate": 1.5},
	}

	runParserTest(t, v, p, obj, 2)

	task["Retry"] = []interface{}{
		map[string]interface{}{"ErrorEquals": 1},
		map[string]interface{}{"ErrorEquals": true},
	}

	runParserTest(t, v, p, obj, 2)

	task["Retry"] = []interface{}{
		map[string]interface{}{"ErrorEquals": []interface{}{1}},
		map[string]interface{}{"ErrorEquals": []interface{}{true}},
	}

	runParserTest(t, v, p, obj, 2)

	task["Retry"] = []interface{}{
		map[string]interface{}{"ErrorEquals": []interface{}{"foo"}},
		map[string]interface{}{"ErrorEquals": []interface{}{"bar"}},
	}

	runParserTest(t, v, p, obj, 0)

	rt := map[string]interface{}{
		"ErrorEquals":     []interface{}{"foo"},
		"IntervalSeconds": "bar",
		"MaxAttempts":     true,
		"BackoffRate":     make(map[string]interface{}),
	}
	task["Retry"] = []interface{}{rt}

	runParserTest(t, v, p, obj, 3)

	rt["IntervalSeconds"] = 0
	rt["MaxAttempts"] = -1
	rt["BackoffRate"] = 0.9

	runParserTest(t, v, p, obj, 3)

	rt["IntervalSeconds"] = 5
	rt["MaxAttempts"] = 99999999
	rt["BackoffRate"] = 1.1

	runParserTest(t, v, p, obj, 1)

	rt["MaxAttempts"] = 99999998

	runParserTest(t, v, p, obj, 0)

	catch := map[string]interface{}{
		"ErrorEquals": []interface{}{"foo"},
		"Next":        "n",
	}
	task["Catch"] = []interface{}{catch}

	runParserTest(t, v, p, obj, 0)

	delete(catch, "Next")

	runParserTest(t, v, p, obj, 1)

	catch["Next"] = true

	runParserTest(t, v, p, obj, 1)

	catch["Next"] = "t"
	delete(catch, "ErrorEquals")

	runParserTest(t, v, p, obj, 1)

	catch["ErrorEquals"] = []interface{}{}

	runParserTest(t, v, p, obj, 1)

	catch["ErrorEquals"] = []interface{}{3}

	runParserTest(t, v, p, obj, 1)

	catch["ErrorEquals"] = []interface{}{"x"}

	runParserTest(t, v, p, obj, 0)

	choice := map[string]interface{}{
		"Type":    "Choice",
		"Default": "x",
	}

	choicesArray := []interface{}{
		map[string]interface{}{
			"Next":           "z",
			"Variable":       "$.a.b",
			"StringLessThan": "xx",
		},
	}
	choice["Choices"] = choicesArray

	delete(states, "Task")
	delete(states, "Fail")

	obj["States"] = states
	states["Choice"] = choice

	runParserTest(t, v, p, obj, 0)

	choice["Next"] = "a"

	runParserTest(t, v, p, obj, 1)

	delete(choice, "Next")

	choice["End"] = true

	runParserTest(t, v, p, obj, 1)

	delete(choice, "End")

	choicesArray = []interface{}{}
	choice["Choices"] = choicesArray

	runParserTest(t, v, p, obj, 1)

	choicesArray = []interface{}{1, "2"}
	choice["Choices"] = choicesArray

	runParserTest(t, v, p, obj, 2)

	choicesArray = []interface{}{
		map[string]interface{}{
			"Next":          "y",
			"Variable":      "$.c.d",
			"NumericEquals": 5,
		},
	}
	choice["Choices"] = choicesArray

	runParserTest(t, v, p, obj, 0)

	nester := map[string]interface{}{
		"And": "foo",
	}

	choicesArray = []interface{}{nester}
	choice["Choices"] = choicesArray

	runParserTest(t, v, p, obj, 2)

	nester["Next"] = "x"

	runParserTest(t, v, p, obj, 1)

	nester["And"] = []interface{}{}

	runParserTest(t, v, p, obj, 1)

	nester["And"] = []interface{}{
		map[string]interface{}{
			"Variable":       "$.a.b",
			"StringLessThan": "xx",
		},
		map[string]interface{}{
			"Variable":      "$.c.d",
			"NumericEquals": 12,
		},
		map[string]interface{}{
			"Variable":      "$.e.f",
			"BooleanEquals": false,
		},
	}

	runParserTest(t, v, p, obj, 0)

	// data types
	bad := []interface{}{
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringEquals": 1},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringLessThan": true},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringGreaterThan": 11.5},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringLessThanEquals": 0},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringGreaterThanEquals": false},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericEquals": "a"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericLessThan": true},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericGreaterThan": []interface{}{3, 4}},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericLessThanEquals": map[string]interface{}{}},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericGreaterThanEquals": "bar"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "BooleanEquals": 3},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampEquals": "a"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampLessThan": 3},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampGreaterThan": true},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampLessThanEquals": false},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampGreaterThanEquals": 3},
	}
	good := []interface{}{
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringEquals": "foo"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringLessThan": "bar"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringGreaterThan": "baz"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringLessThanEquals": "foo"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "StringGreaterThanEquals": "bar"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericEquals": 11},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericLessThan": 12},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericGreaterThan": 13},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericLessThanEquals": 14.3},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "NumericGreaterThanEquals": 14.4},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "BooleanEquals": false},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampEquals": "2016-03-14T01:59:00Z"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampLessThan": "2016-03-14T01:59:00Z"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampGreaterThan": "2016-03-14T01:59:00Z"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampLessThanEquals": "2016-03-14T01:59:00Z"},
		map[string]interface{}{"Variable": "$.a", "Next": "b", "TimestampGreaterThanEquals": "2016-03-14T01:59:00Z"},
	}

	for _, comp := range bad {
		choicesArray = []interface{}{comp}
		choice["Choices"] = choicesArray

		runParserTest(t, v, p, obj, 1)
	}

	for _, comp := range good {
		choicesArray = []interface{}{comp}
		choice["Choices"] = choicesArray

		runParserTest(t, v, p, obj, 0)
	}

	// nesting
	choice["Choices"] = []interface{}{
		map[string]interface{}{
			"Not":  map[string]interface{}{"Variable": "$.type", "StringEquals": "Private"},
			"Next": "Public",
		},
		map[string]interface{}{
			"And": []interface{}{
				map[string]interface{}{"Variable": "$.type", "NumericGreaterThanEquals": 20},
				map[string]interface{}{"Variable": "$.type", "NumericLessThan": 30},
			},
			"Next": "ValueInTwenties",
		},
	}

	runParserTest(t, v, p, obj, 0)

	choice["Choices"] = []interface{}{
		map[string]interface{}{
			"Not":  map[string]interface{}{"Variable": false, "StringEquals": "Private"},
			"Next": "Public",
		},
	}

	runParserTest(t, v, p, obj, 1)

	choice["Choices"] = []interface{}{
		map[string]interface{}{
			"And": []interface{}{
				map[string]interface{}{
					"Variable":                 "$.value",
					"NumericGreaterThanEquals": 20,
					"And":                      true,
				},
				map[string]interface{}{"Variable": "$.value", "NumericLessThan": 44},
			},
			"Next": "ValueInTwenties",
		},
	}

	runParserTest(t, v, p, obj, 2)

	delete(states, "Choice")

	// wait state
	wait := map[string]interface{}{
		"Type":    "Wait",
		"Next":    "z",
		"Seconds": 5,
	}
	states["Wait"] = wait

	runParserTest(t, v, p, obj, 0)

	wait["Seconds"] = "t"

	runParserTest(t, v, p, obj, 1)

	delete(wait, "Seconds")

	wait["SecondsPath"] = 12

	runParserTest(t, v, p, obj, 1)

	delete(wait, "SecondsPath")

	wait["Timestamp"] = false

	runParserTest(t, v, p, obj, 1)

	delete(wait, "Timestamp")

	wait["TimestampPath"] = 33

	runParserTest(t, v, p, obj, 1)

	delete(wait, "TimestampPath")

	wait["Timestamp"] = "2016-03-14T01:59:00Z"

	runParserTest(t, v, p, obj, 0)

	wait = map[string]interface{}{
		"Type":        "Wait",
		"Next":        "z",
		"Seconds":     5,
		"SecondsPath": "$.x",
	}
	states["Wait"] = wait

	runParserTest(t, v, p, obj, 1)

	delete(states, "Wait")

	branches := []interface{}{
		map[string]interface{}{
			"StartAt": "p1",
			"States": map[string]interface{}{
				"p1": map[string]interface{}{
					"Type": "Pass",
					"End":  true,
				},
			},
		},
	}
	para := map[string]interface{}{
		"Type":     "Parallel",
		"End":      true,
		"Branches": branches,
	}
	states["Parallel"] = para

	runParserTest(t, v, p, obj, 0)

	(branches[0].(map[string]interface{}))["StartAt"] = true

	runParserTest(t, v, p, obj, 1)

	para["Branches"] = 3

	runParserTest(t, v, p, obj, 1)

	para["Branches"] = []interface{}{}

	runParserTest(t, v, p, obj, 0)

	para["Branches"] = []interface{}{map[string]interface{}{}}

	runParserTest(t, v, p, obj, 2)

	para["Branches"] = []interface{}{
		map[string]interface{}{
			"StartAt": "p1",
			"States": map[string]interface{}{
				"p1": map[string]interface{}{
					"Type": "foo",
					"End":  true,
				},
			},
		},
	}

	runParserTest(t, v, p, obj, 2)

	para["Branches"] = []interface{}{
		map[string]interface{}{
			"foo":     1,
			"StartAt": "p1",
			"States": map[string]interface{}{
				"p1": map[string]interface{}{
					"Type": "Pass",
					"End":  true,
				},
			},
		},
	}

	runParserTest(t, v, p, obj, 1)
}

func runParserTest(t *testing.T, v *NodeValidator, p *J2119Parser, jsonObject interface{}, wantedErrorCount int) {
	t.Helper()

	problems := NewProblems()
	node := *NewNode(jsonObject)

	v.Validate(node, p.root, []string{p.root}, problems)

	if wantedErrorCount != -1 {
		assert.Equal(t, wantedErrorCount, problems.Len())
	}
}
