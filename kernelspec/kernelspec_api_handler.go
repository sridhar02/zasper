package kernelspec

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// Handler to serve kernel resources (like logos)
func ServeKernelResource(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	kernelName := vars["kernel"]
	resourcePath := vars["resource"]

	fullPath := getResourceFile(kernelName, resourcePath)

	resourceData, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	ext := strings.ToLower(strings.TrimPrefix(resourcePath, "."))
	var contentType string
	switch ext {
	case "png":
		contentType = "image/png"
	case "svg":
		contentType = "image/svg+xml"
	default:
		contentType = "application/octet-stream" // default for unknown file types
	}

	// Set the content type for the response
	w.Header().Set("Content-Type", contentType)
	// Serve the resource data
	w.Write(resourceData)
}

func SingleKernelspecAPIHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Get request received")

	vars := mux.Vars(req)
	kernelName := vars["kernelName"]
	log.Println("kernelName :", kernelName)

	kspec := GetKernelSpec(kernelName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kspec)
}

func KernelspecAPIHandler(w http.ResponseWriter, req *http.Request) {
	response := KernelspecResponse{
		Default:    "python3",
		Kernespecs: make(map[string]KernelspecModel),
	}

	available_kernelspec := GetAllSpecs()

	for kernelName, kernelInfo := range available_kernelspec {
		response.Kernespecs[kernelName] = KernelspecModel{
			Name:      kernelName,
			Spec:      kernelInfo.Spec,
			Resources: getResources(kernelName, kernelInfo.ResourceDir),
		}

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
