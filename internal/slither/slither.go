package slither

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunAll(projectRoot string, scopeFiles []string) ([]*SlitherResult, error) {
	tmpFile := filepath.Join(projectRoot, ".slither_out.json")
	defer os.Remove(tmpFile)

	fmt.Println("       Running slither on project...")
	cmd := exec.Command("slither", ".", "--json", tmpFile)
	cmd.Dir = projectRoot
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("       slither exited with: %v (checking output anyway...)\n", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		return nil, fmt.Errorf("slither output not found: %v", err)
	}

	var raw slitherRaw
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse slither JSON: %v", err)
	}

	return filterToScope(raw, scopeFiles), nil
}

func filterToScope(raw slitherRaw, scopeFiles []string) []*SlitherResult {
	inScope := make(map[string]bool)
	for _, f := range scopeFiles {
		inScope[filepath.ToSlash(f)] = true
	}

	fileMap := make(map[string]*SlitherResult)

	matchScope := func(slitherPath string) string {
		normalized := filepath.ToSlash(slitherPath)
		for scopeFile := range inScope {
			if strings.HasSuffix(normalized, scopeFile) {
				return scopeFile
			}
		}
		return ""
	}

	getResult := func(slitherPath string) *SlitherResult {
		key := matchScope(slitherPath)
		if key == "" {
			return nil
		}
		if _, exists := fileMap[key]; !exists {
			fileMap[key] = &SlitherResult{
				FilePath: key,
				Success:  raw.Success,
			}
		}
		return fileMap[key]
	}

	// Map contracts
	for _, c := range raw.Results.Contracts {
		r := getResult(c.SourceMapping.Filename)
		if r == nil {
			continue
		}
		contract := Contract{Name: c.Name}
		for _, inh := range c.Inheritance {
			contract.Inheritance = append(contract.Inheritance, inh.Name)
		}
		for _, f := range c.Functions {
			fn := Function{Name: f.Name, Visibility: f.Visibility}
			for _, m := range f.Modifiers {
				fn.Modifiers = append(fn.Modifiers, m.Name)
			}
			for _, p := range f.Parameters {
				fn.Parameters = append(fn.Parameters, p.Type)
			}
			for _, ret := range f.ReturnType {
				fn.Returns = append(fn.Returns, ret.Type)
			}
			contract.Functions = append(contract.Functions, fn)
		}
		for _, sv := range c.StateVariables {
			contract.StateVars = append(contract.StateVars, StateVar{
				Name: sv.Name, Type: sv.Type, Visibility: sv.Visibility,
			})
		}
		r.Contracts = append(r.Contracts, contract)
	}

	// Map detectors
	for _, d := range raw.Results.Detectors {
		for _, elem := range d.Elements {
			r := getResult(elem.SourceMapping.Filename)
			if r == nil {
				continue
			}
			r.Detectors = append(r.Detectors, Detector{
				Check:       d.Check,
				Impact:      d.Impact,
				Confidence:  d.Confidence,
				Description: d.Description,
			})
			break
		}
	}

	var results []*SlitherResult
	for _, r := range fileMap {
		results = append(results, r)
	}
	return results
}
