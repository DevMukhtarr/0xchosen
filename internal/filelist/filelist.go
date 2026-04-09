package filelist

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Generate(srcDir, outputFile string) error {
	re := regexp.MustCompile(`function\s+[^\(]+\([^\)]*\)\s+.*\b(public|external)\b`)

	var entries []string

	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".sol") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for _, line := range strings.Split(string(content), "\n") {
			if re.MatchString(line) {
				// Format same as grep output: "filepath:line"
				entries = append(entries, filepath.ToSlash(path)+":"+strings.TrimSpace(line))
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}
