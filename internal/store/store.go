package store

import (
	"encoding/json"
	"os"

	"github.com/devmukhtarr/0xchosen/internal/extractor"
)

func Save(facts *extractor.Facts, outputFile string) error {
	data, err := json.MarshalIndent(facts, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outputFile, data, 0644)
}
