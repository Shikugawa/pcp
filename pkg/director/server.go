package director

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Shikugawa/pcp/pkg/filter"
	"github.com/Shikugawa/pcp/pkg/manager"
	"github.com/Shikugawa/pcp/pkg/nodes"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
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
	for _, node := range nodes.ManagedNodes.GetAll() {
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, reader, err := r.FormFile("wasm-code")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wasmCode := make([]byte, reader.Size)
	file.Read(wasmCode)

	specifier, err := filter.ParseFileName(reader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.envoyFilterManager.Storage.Add(*specifier, wasmCode); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
		log.Println("Admin server started...")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln("Server closed with error:", err)
		}
	}()

	return srv
}
