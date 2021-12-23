package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func GetJSONFromLocalFile(path string) (interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can not open local file \"%s\": %w", path, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("can not read local file \"%s\": %w", path, err)
	}

	var j interface{}

	err = json.Unmarshal(data, &j)
	if err != nil {
		return nil, fmt.Errorf("can not unmarshal local file \"%s\": %w", path, err)
	}

	return j, nil
}
