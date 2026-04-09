package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func ExtractFilePaths(funcsFile string) ([]string, error) {
	data, err := os.ReadFile(funcsFile)
	if err != nil {
		return nil, err
	}

	var entries []string
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var files []string
	for _, entry := range entries {
		parts := strings.SplitN(entry, ":", 2)
		if len(parts) < 1 {
			continue
		}
		filePath := strings.TrimSpace(parts[0])
		filePath = strings.Trim(filePath, `"`)
		filePath = filepath.ToSlash(filePath)
		if filePath != "" && !seen[filePath] {
			seen[filePath] = true
			files = append(files, filePath)
		}
	}

	return files, nil
}
