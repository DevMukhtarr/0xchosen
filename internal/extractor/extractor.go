package extractor

import (
	"github.com/devmukhtarr/0xchosen/internal/dependency"
	"github.com/devmukhtarr/0xchosen/internal/slither"
)

type Facts struct {
	Files           []FileFacts
	DependencyGraph *dependency.Graph
}

type FileFacts struct {
	FilePath  string
	Contracts []ContractFacts
	Detectors []slither.Detector
}

type ContractFacts struct {
	Name          string
	Inheritance   []string
	EntryPoints   []slither.Function
	AccessControl []AccessControlFact
	StateVars     []slither.StateVar
	Events        []string
}

type AccessControlFact struct {
	FunctionName string
	Modifiers    []string
}

func Extract(results []*slither.SlitherResult, graph *dependency.Graph) *Facts {
	facts := &Facts{DependencyGraph: graph}

	for _, result := range results {
		if !result.Success && result.Error != "" {
			continue
		}

		fileFact := FileFacts{
			FilePath:  result.FilePath,
			Detectors: result.Detectors,
		}

		for _, contract := range result.Contracts {
			cf := ContractFacts{
				Name:        contract.Name,
				Inheritance: contract.Inheritance,
				StateVars:   contract.StateVars,
			}

			for _, fn := range contract.Functions {
				// Only extract public/external entry points
				if fn.Visibility == "public" || fn.Visibility == "external" {
					cf.EntryPoints = append(cf.EntryPoints, fn)

					// Track access control
					if len(fn.Modifiers) > 0 {
						cf.AccessControl = append(cf.AccessControl, AccessControlFact{
							FunctionName: fn.Name,
							Modifiers:    fn.Modifiers,
						})
					}
				}
			}

			fileFact.Contracts = append(fileFact.Contracts, cf)
		}

		facts.Files = append(facts.Files, fileFact)
	}

	return facts
}
