package j2119

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetValidator(t *testing.T) *Validator {
	t.Helper()

	validator, err := NewValidator("../data/AWL.j2119")
	if err != nil {
		t.Fatal("can not create validator")
	}

	return validator
}

func TestValidator_Good(t *testing.T) {
	t.Parallel()

	v := GetValidator(t)
	j := GetJSONObjectFromFile(t, "good.json")
	p := v.ValidateJSONStruct(j)

	assert.Equal(t, 0, p.Len())
}

func TestValidator_WithArrayResult(t *testing.T) {
	t.Parallel()

	v := GetValidator(t)
	j := GetJSONObjectFromFile(t, "withArrayResult.json")
	p := v.ValidateJSONStruct(j)

	assert.Equal(t, 0, p.Len())
}

func TestValidator_WithObjectResult(t *testing.T) {
	t.Parallel()

	v := GetValidator(t)
	j := GetJSONObjectFromFile(t, "withObjectResult.json")
	p := v.ValidateJSONStruct(j)

	assert.Equal(t, 0, p.Len())
}
