package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func ResolveLatestVersion(module string) (string, error) {
	url := fmt.Sprintf("https://proxy.golang.org/%s/@latest", module)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to resolve latest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to resolve latest: status %d", resp.StatusCode)
	}

	var data struct {
		Version string
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", fmt.Errorf("failed to parse version json: %w", err)
	}

	return data.Version, nil
}

func CompareVersions(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	s1 := strings.Split(v1, ".")
	s2 := strings.Split(v2, ".")

	for i := range 3 {
		var a, b int
		if i < len(s1) {
			fmt.Sscanf(s1[i], "%d", &a)
		}
		if i < len(s2) {
			fmt.Sscanf(s2[i], "%d", &b)
		}
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
	}
	return 0
}
