package filter

import (
	"fmt"
	"strings"
)

type FilterSpecifier struct {
	FilterType string `json:"filter_type"`
	FilterName string `json:"filter_name"`
}

func (f *FilterSpecifier) String() string {
	return fmt.Sprintf("%s.%s", f.FilterType, f.FilterName)
}

func StringToSpecifier(s string) *FilterSpecifier {
	splitedSpecifier := strings.Split(s, ".")
	if len(splitedSpecifier) != 2 {
		return nil
	}
	return &FilterSpecifier{FilterType: splitedSpecifier[0], FilterName: splitedSpecifier[1]}
}
