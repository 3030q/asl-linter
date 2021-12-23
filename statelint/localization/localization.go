package localization

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	path2 "path"
	"sync"
)

const (
	DefaultLanguage           = "en"
	DefaultLocalizationFolder = "./langs"
)

var (
	once      sync.Once
	singleton Localizer
)

type Localizer interface {
	// SetLocalization load language file. By default, language is en.
	SetLocalization(lang string) error
	// GetString returns string by given name
	GetString(name string) string
}

type localizer struct {
	currentLanguage string
	folder          string
	jsonFile        map[string]string
	mu              sync.RWMutex
}

func GetLocalizerFromFile(localizationFolder string) (Localizer, error) {
	var err error

	once.Do(func() {
		singleton = &localizer{
			currentLanguage: "",
			jsonFile:        nil,
			folder:          localizationFolder,
			mu:              sync.RWMutex{},
		}
		err = singleton.SetLocalization(DefaultLanguage)
	})

	if err != nil {
		return nil, err
	}

	return singleton, nil
}

func GetLocalizer() (Localizer, error) {
	return GetLocalizerFromFile(DefaultLocalizationFolder)
}

func GetLocalizerOrPanic() Localizer {
	l, err := GetLocalizer()
	if err != nil {
		panic("can not get localizer")
	}

	return l
}

func (l *localizer) SetLocalization(lang string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.currentLanguage = lang

	if err := l.loadLocalization(); err != nil {
		return err
	}

	return nil
}

func (l *localizer) GetString(name string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if str, ok := l.jsonFile[name]; ok {
		return str
	}

	panic(fmt.Sprintf("can not find string %s in %s.json file", name, l.currentLanguage))
}

func (l *localizer) loadLocalization() error {
	pathToFile := path2.Join(l.folder, l.currentLanguage+".json")

	exists, err := fileExists(pathToFile)
	if err != nil {
		return err
	}

	if !exists {
		return errors.New("can not find localization file")
	}

	data, err := os.ReadFile(pathToFile)
	if err != nil {
		return fmt.Errorf("can not read localization file %s: %w", pathToFile, err)
	}

	var jsonData map[string]string
	err = json.Unmarshal(data, &jsonData)

	if err != nil {
		return fmt.Errorf("error when unmarshal json: %w", err)
	}

	l.jsonFile = jsonData

	return nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, fmt.Errorf("error when check filepath existence %s, err: %w", path, err)
}
