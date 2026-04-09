package relationship

import (
	"fmt"
	"strings"

	"github.com/devmukhtarr/0xchosen/internal/extractor"
)

type Summary struct {
	InheritanceTree  []string
	CallSurface      []string
	AdminFunctions   []string
	KeyStateVars     []string
	DetectorWarnings []string
}

func Build(facts *extractor.Facts) *Summary {
	summary := &Summary{}

	for _, file := range facts.Files {
		for _, contract := range file.Contracts {
			// Inheritance tree
			if len(contract.Inheritance) > 0 {
				summary.InheritanceTree = append(summary.InheritanceTree,
					fmt.Sprintf("%s inherits: %s",
						contract.Name,
						strings.Join(contract.Inheritance, ", ")))
			}

			// Entry points
			for _, fn := range contract.EntryPoints {
				summary.CallSurface = append(summary.CallSurface,
					fmt.Sprintf("%s.%s() [%s]",
						contract.Name, fn.Name, fn.Visibility))
			}

			// Admin/access controlled functions
			for _, ac := range contract.AccessControl {
				summary.AdminFunctions = append(summary.AdminFunctions,
					fmt.Sprintf("%s.%s() — modifiers: %s",
						contract.Name,
						ac.FunctionName,
						strings.Join(ac.Modifiers, ", ")))
			}

			// Key state vars (mappings and balances)
			for _, sv := range contract.StateVars {
				if strings.Contains(sv.Type, "mapping") ||
					strings.Contains(strings.ToLower(sv.Name), "balance") ||
					strings.Contains(strings.ToLower(sv.Name), "owner") ||
					strings.Contains(strings.ToLower(sv.Name), "admin") {
					summary.KeyStateVars = append(summary.KeyStateVars,
						fmt.Sprintf("%s: %s %s", sv.Name, sv.Type, sv.Visibility))
				}
			}
		}

		// High/medium impact detector warnings
		for _, d := range file.Detectors {
			if d.Impact == "High" || d.Impact == "Medium" {
				summary.DetectorWarnings = append(summary.DetectorWarnings,
					fmt.Sprintf("[%s][%s] %s", d.Impact, d.Check, d.Description))
			}
		}
	}

	return summary
}
