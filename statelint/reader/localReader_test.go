package reader

import "testing"

func TestLocalReader_RightPath(t *testing.T) {
	t.Parallel()

	base := "../testdata/"
	testCases := []string{
		base + "good.json",
		base + "noTerminal.json",
	}

	for _, testCase := range testCases {
		_, err := GetJSONFromLocalFile(testCase)
		if err != nil {
			t.Fatalf("should be nil err, but err = %s", err.Error())
		}
	}
}

func TestLocalReader_NotExistedPaths(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"a",
		"b",
	}

	for _, testCase := range testCases {
		_, err := GetJSONFromLocalFile(testCase)
		if err == nil {
			t.Fatal("should be not nill err, but err = nil")
		}
	}
}
