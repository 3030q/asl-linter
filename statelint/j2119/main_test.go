package j2119

import (
	"log"
	"os"
	"statelint/localization"
	"testing"
)

// Test main needed to initialize Localizer to testing case with problems
func TestMain(m *testing.M) {
	// set up localizer with default language
	_, err := localization.GetLocalizerFromFile("../langs")
	if err != nil {
		log.Fatal("can not init localizer")
	}

	code := m.Run()

	os.Exit(code)
}
