package filter

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wasmRegistoryPath = "/wasm"

type FilterStorage struct {
	filters map[string]string
	mux     sync.Mutex
}

func NewFilterStorage(path string) *FilterStorage {
	wasmRegistoryPath = path

	storage := &FilterStorage{
		filters: map[string]string{},
	}
	if err := filepath.Walk(wasmRegistoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return filepath.SkipDir
		}

		splitedPath := strings.Split(path, "/")
		// Splited path is expected to have only "wasm" and <filter_type>
		if len(splitedPath) != 2 {
			return filepath.SkipDir
		}
		filterType := splitedPath[1]

		splitedFileName := strings.Split(info.Name(), ".")
		// filename is expected to be wasm code
		if splitedFileName[len(splitedFileName)-1] != "wasm" {
			return filepath.SkipDir
		}
		filterName := strings.Join(splitedFileName[0:len(splitedFileName)-1], ".")

		specifier := FilterSpecifier{
			FilterType: filterType,
			FilterName: filterName,
		}
		storage.filters[specifier.String()] = WasmCodePath(&specifier)

		return nil
	}); err != nil {
		return nil
	}
	return storage
}

func (f *FilterStorage) GetAll() []*FilterSpecifier {
	var specifiers []*FilterSpecifier
	for dir, _ := range f.filters {
		if spec := wasmCodeDirToSpecifier(dir); spec != nil {
			specifiers = append(specifiers, spec)
		}
	}
	return specifiers
}

func (f *FilterStorage) Add(filter FilterSpecifier, wasmCode []byte) error {
	if _, err := os.Stat(wasmCodeDir(&filter)); os.IsNotExist(err) {
		if err = os.Mkdir(wasmCodeDir(&filter), 0644); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(
		WasmCodePath(&filter), wasmCode, 0644); err != nil {
		return err
	}

	f.mux.Lock()
	defer f.mux.Unlock()
	f.filters[filter.String()] = WasmCodePath(&filter)

	return nil
}

func (f *FilterStorage) Del(filter FilterSpecifier) error {
	if !f.ExistFilter(filter) {
		return errors.New("Filter not found")
	}

	if err := os.Remove(WasmCodePath(&filter)); err != nil {
		return errors.New("Unable to remove file")
	}

	f.mux.Lock()
	defer f.mux.Unlock()
	delete(f.filters, filter.String())

	return nil
}

func (f *FilterStorage) ExistFilter(filter FilterSpecifier) bool {
	for k, _ := range f.filters {
		if k == filter.String() {
			return true
		}
	}
	return false
}
