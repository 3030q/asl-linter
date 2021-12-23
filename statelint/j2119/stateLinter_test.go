package j2119

import (
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetJSONObjectFromFile(t *testing.T, filename string) interface{} {
	t.Helper()

	p := "../testdata/" + filename

	f, err := os.Open(p)
	if err != nil {
		t.Fatalf("Can not open file \"%s\": %s", p, err)
	}
	defer f.Close()

	var j interface{}

	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Can not read file \"%s\": %s", p, err)
	}

	err = json.Unmarshal(data, &j)
	if err != nil {
		t.Fatalf("Can not unmarshal json file \"%s\": %s", p, err)
	}

	return j
}

func TestStateLinter_JSONRightValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		filename      string
		problemsCount int
	}{
		{"minimalFailState.json", 0},
		{"emptyErrorEqualsOnCatch.json", 1},
		{"emptyErrorEqualsOnRetry.json", 1},
		{"passWithParameters.json", 0},
		{"taskWithParameters.json", 0},
		{"choiceWithParameters.json", 1},
		{"waitWithParameters.json", 1},
		{"succeedWithParameters.json", 1},
		{"failWithParameters.json", 1},
		{"parallelWithParameters.json", 0},
		{"mapWithParameters.json", 0},
		{"parameterPathProblems.json", 5},
		{"passWithResultpath.json", 0},
		{"taskWithResultpath.json", 0},
		{"choiceWithResultpath.json", 1},
		{"waitWithResultpath.json", 1},
		{"succeedWithResultpath.json", 1},
		{"failWithResultpath.json", 1},
		{"parallelWithResultpath.json", 0},
		{"passWithIoPathContextObject.json", 0},
		{"choiceWithContextObject.json", 0},
		{"mapWithItemspathContextObject.json", 0},
		{"taskWithDynamicTimeouts.json", 0},
		{"passWithNullInputpath.json", 0},
		{"passWithNullOutputpath.json", 0},
		{"mapWithNullItemspath.json", 1},
		{"taskWithResultselector.json", 0},
		{"parallelWithResultselector.json", 0},
		{"mapWithResultselector.json", 0},
		{"passWithResultselector.json", 1},
		{"waitWithResultselector.json", 1},
		{"failWithResultselector.json", 1},
		{"succeedWithResultselector.json", 1},
		{"choiceWithResultselector.json", 1},
		{"statesArrayInvocation.json", 0},
		{"statesFormatInvocation.json", 0},
		{"statesStringtojsonInvocation.json", 0},
		{"statesJsontostringInvocation.json", 0},
		{"statesArrayInvocationLeftpad.json", 1},
		{"invalidFunctionInvocation.json", 1},
		{"passWithIntrinsicFunctionInputpath.json", 1},
		{"taskWithStaticAndDynamicTimeout.json", 1},
		{"taskWithStaticAndDynamicHeartbeat.json", 1},
	}

	linter, err := NewStateLinterFromFile("." + DefaultStateMachinePath)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	for _, testCase := range testCases {
		jsonObject := GetJSONObjectFromFile(t, testCase.filename)
		problems := linter.ValidateJSONStruct(jsonObject)

		assert.Equalf(
			t,
			testCase.problemsCount,
			problems.Len(),
			"wrong problems count with file \"%s\"",
			testCase.filename,
		)
	}
}
