package filter

import (
	"errors"
	"fmt"
	"strings"
)

func wasmCodePath(filter FilterSpecifier) string {
	return wasmRegistoryPath + "/" + filter.FilterType + "/" + filter.FilterName + ".wasm"
}

func wasmCodeDir(filter FilterSpecifier) string {
	return wasmRegistoryPath + "/" + filter.FilterType
}

func ParseFileName(fileName string) (*FilterSpecifier, error) {
	splitedFilename := strings.Split(fileName, ".")
	if len(splitedFilename) != 3 || splitedFilename[2] != "wasm" {
		return nil, errors.New(fmt.Sprintf("Invalid format of file name: %s", fileName))
	}
	return &FilterSpecifier{FilterType: splitedFilename[0], FilterName: splitedFilename[1]}, nil
}
