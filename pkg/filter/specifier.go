package filter

import "fmt"

type FilterSpecifier struct {
	FilterType string `json:"filter_type"`
	FilterName string `json:"filter_name"`
}

func (f *FilterSpecifier) String() string {
	return fmt.Sprintf("%s.%s", f.FilterType, f.FilterName)
}
