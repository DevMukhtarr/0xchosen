package dependency

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Graph struct {
	// file -> list of files it imports
	Imports map[string][]string
	// file -> list of files that import it
	ImportedBy map[string][]string
}

var importRegex = regexp.MustCompile(`^\s*import\s+["']([^"']+)["']`)
var importBracesRegex = regexp.MustCompile(`from\s+["']([^"']+)["']`)

func Build(files []string) (*Graph, error) {
	graph := &Graph{
		Imports:    make(map[string][]string),
		ImportedBy: make(map[string][]string),
	}

	for _, file := range files {
		imports, err := parseImports(file)
		if err != nil {
			// skip files we can't read
			continue
		}
		graph.Imports[file] = imports
		for _, imp := range imports {
			graph.ImportedBy[imp] = append(graph.ImportedBy[imp], file)
		}
	}

	return graph, nil
}

func parseImports(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var imports []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		// Match: import "path" or import './path'
		if m := importRegex.FindStringSubmatch(line); len(m) > 1 {
			imports = append(imports, normalizeImport(filePath, m[1]))
			continue
		}

		// Match: import {X} from "path"
		if m := importBracesRegex.FindStringSubmatch(line); len(m) > 1 {
			imports = append(imports, normalizeImport(filePath, m[1]))
		}
	}

	return imports, scanner.Err()
}

func normalizeImport(sourceFile, importPath string) string {
	// Leave package imports as-is (@openzeppelin, etc.)
	if strings.HasPrefix(importPath, "@") {
		return importPath
	}

	// Resolve relative to the source file's directory
	sourceDir := filepath.Dir(sourceFile)
	resolved := filepath.Join(sourceDir, importPath)
	return filepath.ToSlash(filepath.Clean(resolved))
}
