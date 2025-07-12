package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type proxyOrigin struct {
	Hash string `json:"Hash"`
}

type proxyMeta struct {
	Version string      `json:"Version"`
	Time    time.Time   `json:"Time"`
	Origin  proxyOrigin `json:"Origin"`
}

type ModuleMetadata struct {
	Version string    `json:"Version"`
	Time    time.Time `json:"Time"`
	Hash    string    `json:"Hash,omitempty"`
}

func FetchModuleMetadata(module, version string) (*ModuleMetadata, error) {
	url := fmt.Sprintf("https://proxy.golang.org/%s/@v/%s.info", module, version)
	if version == "latest" {
		url = fmt.Sprintf("https://proxy.golang.org/%s/@latest", module)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to resolve version %q for %s: status %d", version, module, resp.StatusCode)
	}

	var data proxyMeta
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &ModuleMetadata{
		Version: data.Version,
		Time:    data.Time,
		Hash:    data.Origin.Hash,
	}, nil
}
