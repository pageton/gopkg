package core

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func ScanImports(root string) ([]string, error) {
	var imports []string
	seen := map[string]bool{}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return nil
		}
		for _, imp := range node.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.Contains(path, ".") && !seen[path] {
				seen[path] = true
				imports = append(imports, path)
			}
		}
		return nil
	})

	return imports, err
}
