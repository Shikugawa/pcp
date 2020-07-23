package filter

import "fmt"

type FilterSpecifier struct {
	FilterType string
	FilterName string
}

func (f *FilterSpecifier) String() string {
	return fmt.Sprintf("%s.%s", f.FilterType, f.FilterName)
}
