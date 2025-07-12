package core

import (
	"fmt"
	"os/exec"
	"strings"
)

func AddReplaceToGoMod(module, localPath, version string) error {
	module = strings.Split(module, "@")[0]

	_ = exec.Command("go", "mod", "edit", "-require", fmt.Sprintf("%s@%s", module, version)).Run()

	cmd := exec.Command("go", "mod", "edit", "-replace", fmt.Sprintf("%s=%s", module, localPath))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add replace for %s: %w", module, err)
	}

	return nil
}
