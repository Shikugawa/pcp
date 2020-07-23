package director

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Shikugawa/pcp/filter"

	node "github.com/Shikugawa/pcp/nodes"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

	"github.com/Shikugawa/pcp/manager"
)

type Server struct {
	envoyFilterManager *manager.EnvoyFilterManager
}

func NewServer(envoyFilterManager *manager.EnvoyFilterManager) *Server {
	return &Server{
		envoyFilterManager: envoyFilterManager,
	}
}

func (s *Server) enableFilter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request UpdateFilterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var targetNodes []*core.Node
	for _, nodes := range request.Nodes {
		slicedNodes := strings.Split(nodes, "/")
		if len(slicedNodes) != 2 {
			continue
		}

		targetNodes = append(targetNodes, &core.Node{
			Cluster: slicedNodes[0],
			Id:      slicedNodes[1],
		})
	}

	if err := s.envoyFilterManager.Append(request.FilterType, request.FilterName, targetNodes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (s *Server) disableFilter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request UpdateFilterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var targetNodes []*core.Node
	for _, nodes := range request.Nodes {
		slicedNodes := strings.Split(nodes, "/")
		if len(slicedNodes) != 2 {
			continue
		}

		targetNodes = append(targetNodes, &core.Node{
			Cluster: slicedNodes[0],
			Id:      slicedNodes[1],
		})
	}

	if err := s.envoyFilterManager.RemoveFilter(request.FilterType, request.FilterName, targetNodes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (s *Server) nodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var nodesStr []string
	for _, node := range node.Nodes {
		nodesStr = append(nodesStr, node.Cluster+"/"+node.Id)
	}

	resp, err := json.Marshal(&NodesResponse{
		Nodes: nodesStr,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("%s\n", string(resp))))
	return
}

func (s *Server) uploadWasm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/wasm" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var request WasmUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var wasmCode []byte
	_, err := base64.StdEncoding.Decode(wasmCode, []byte(request.Contents))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.envoyFilterManager.StorageDriver.Storage.Add(filter.FilterSpecifier{
		FilterType: request.FilterType,
		FilterName: request.FilterName,
	}, wasmCode)

	w.WriteHeader(http.StatusOK)
	return
}

func (s *Server) Start(port string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/enable", s.enableFilter)
	mux.HandleFunc("/disable", s.disableFilter)
	mux.HandleFunc("/nodes", s.nodes)
	mux.HandleFunc("/upload", s.uploadWasm)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("Server closed with error:", err)
		}
	}()

	return srv
}
