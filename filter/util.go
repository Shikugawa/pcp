package filter

func wasmCodeDir(filter FilterSpecifier) string {
	return wasmRegistoryPath + "/" + filter.FilterType + "/" + filter.FilterName + ".wasm"
}
