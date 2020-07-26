package filter

import (
	"errors"
	"fmt"
	"strings"
)

func WasmCodePath(filter FilterSpecifier) string {
	return wasmRegistoryPath + "/" + filter.FilterType + "/" + filter.FilterName + ".wasm"
}

func wasmCodeDir(filter FilterSpecifier) string {
	return wasmRegistoryPath + "/" + filter.FilterType
}

func wasmCodeDirToSpecifier(dir string) *FilterSpecifier {
	splitedDir := strings.Split(dir, ".")
	if len(splitedDir) != 2 {
		return nil
	}
	return &FilterSpecifier{
		FilterType: splitedDir[0],
		FilterName: splitedDir[1],
	}
}

func ParseFileName(fileName string) (*FilterSpecifier, error) {
	splitedFilename := strings.Split(fileName, ".")
	if len(splitedFilename) != 3 || splitedFilename[2] != "wasm" {
		return nil, errors.New(fmt.Sprintf("Invalid format of file name: %s", fileName))
	}
	return &FilterSpecifier{FilterType: splitedFilename[0], FilterName: splitedFilename[1]}, nil
}
