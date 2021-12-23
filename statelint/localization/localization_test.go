package localization

import (
	"sync"
	"testing"
)

const testDefaulLocalizationFolder = "../langs"

func cleanup() {
	once = sync.Once{}
	singleton = nil
}

// nolint:paralleltest
func TestLocalizer_InitRight(t *testing.T) {
	t.Cleanup(cleanup)

	_, err := GetLocalizerFromFile(testDefaulLocalizationFolder)
	if err != nil {
		t.Fatalf("should init without err, but have: %s", err.Error())
	}
}

// nolint:paralleltest
func TestLocalizer_InitWrong(t *testing.T) {
	t.Cleanup(cleanup)

	if _, err := GetLocalizer(); err == nil {
		t.Fatal("should init with err, but err is nil")
	}
}

// nolint:paralleltest
func TestLocalizer_GetString(t *testing.T) {
	t.Cleanup(cleanup)

	l, err := GetLocalizerFromFile(testDefaulLocalizationFolder)
	if err != nil {
		t.Fatalf("should init without err, but have: %s", err.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Panic occurred!")
		}
	}()

	l.GetString("OnlyOneConstraint")
}

// nolint:paralleltest
func TestLocalizer_GetWrongString(t *testing.T) {
	t.Cleanup(cleanup)

	l, err := GetLocalizerFromFile(testDefaulLocalizationFolder)
	if err != nil {
		t.Fatalf("should init without err, but have: %s", err.Error())
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Panic not occurred!")
		}
	}()

	l.GetString("test")
}

// nolint:paralleltest
func TestLocalizer_GetLocalizerOrPanicReturnPanic(t *testing.T) {
	t.Cleanup(cleanup)

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Panic not occurred!")
		}
	}()

	GetLocalizerOrPanic()
}

// nolint:paralleltest
func TestLocalizer_SetRuLocalization(t *testing.T) {
	t.Cleanup(cleanup)

	l, err := GetLocalizerFromFile(testDefaulLocalizationFolder)
	if err != nil {
		t.Fatalf("should init without err, but have: %s", err.Error())
	}

	err = l.SetLocalization("ru")
	if err != nil {
		t.Fatalf("should not return err, but have: %s", err.Error())
	}
}

// nolint:paralleltest
func TestLocalizer_SetUndefinedLocalization(t *testing.T) {
	t.Cleanup(cleanup)

	l, err := GetLocalizerFromFile(testDefaulLocalizationFolder)
	if err != nil {
		t.Fatalf("should init without err, but have: %s", err.Error())
	}

	err = l.SetLocalization("undefined")
	if err == nil {
		t.Fatal("should return err, but err is nil")
	}
}
