package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	odmPlugin "github.com/hembrow-innovations/odm-plugin" // Replace with odm-plugin library path

	"github.com/hashicorp/go-plugin"
)

// ExecuterImpl is the concrete implementation of our Executer interface.
type ExecuterImpl struct {
	plugin.Plugin
}

type Options struct {
	Action    string   `json:"action"`
	Base      string   `json:"base"`
	Output    string   `json:"output"`
	Overrides []string `json:"overrides"`

	ProjectPath   string   `json:"projectPath"`   // Path to the root of you project
	ProjectFolder string   `json:"projectFolder"` // name of the folder within your projects name
	Projects      []string `json:"projects"`      // an array of string (names of the service "folder name")
	BasePath      string   `json:"basePath"`      // path from your project root to the folder containing a base level file (docker-compose.yml etc)
	ConfigFolder  string   `json:"configFolder"`  // name of folder containing config

}

type ExecutionRequestBody struct {
	Args    map[string]string `json:"args"`
	Options Options           `json:"options"`
	Input   string            `json:"input"`
}

// Greet implements the Greeter interface.
// This signature MUST match shared.Greeter's Greet method.
func (g *ExecuterImpl) Execute(ctx context.Context, body string) (string, error) {

	request := &ExecutionRequestBody{}

	err := json.Unmarshal([]byte(body), request)
	if err != nil {
		return "", err
	}

	if request.Options.Action == "merge" {
		Merge(request)
	} else {
		return "", fmt.Errorf("%s action not found", request.Options.Action)
	}

	log.Printf("Plugin: Execute called with body: \n\tArguments:%s\n\tInput: %s\n\tOptions: %s", request.Args, request.Input, request.Options)
	return "Success", nil
}

func main() {
	log.Println("Starting Tester plugin...")

	// The plugin must export an implementation of the Greeter interface.
	var handshakeConfig = odmPlugin.HandshakeConfig
	var pluginMap = map[string]plugin.Plugin{
		"executer": &odmPlugin.ExecuterPlugin{Impl: &ExecuterImpl{}},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		// GRPCServer:      nil,
		// GRPCProvider:    nil, // We're using standard RPC for simplicity
	})
	log.Println("Greeter plugin finished serving.")
	// This line should never be reached under normal circumstances
	log.Println("Plugin: Serve returned (this shouldn't happen)")
}
