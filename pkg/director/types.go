package director

import "github.com/Shikugawa/pcp/pkg/filter"

type UpdateFilterRequest struct {
	FilterType string   `json:"filter_type"`
	FilterName string   `json:"filter_name"`
	Nodes      []string `json:"nodes"`
}

type NodesResponse struct {
	Nodes []string `json:"nodes"`
}

type WasmUploadRequest struct {
	FilterType string `json:"filter_type"`
	FilterName string `json:"filter_name"`
	Contents   string `json:"contents"`
}

type WasmListResponse struct {
	Filters []filter.FilterSpecifier `json:"filters"`
}
