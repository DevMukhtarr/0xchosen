package slither

type SlitherResult struct {
	FilePath  string
	Contracts []Contract
	Detectors []Detector
	Success   bool
	Error     string
}

type Contract struct {
	Name        string     `json:"name"`
	Inheritance []string   `json:"inheritance"`
	Functions   []Function `json:"functions"`
	StateVars   []StateVar `json:"state_variables"`
}

type Function struct {
	Name       string   `json:"name"`
	Visibility string   `json:"visibility"`
	Modifiers  []string `json:"modifiers"`
	Parameters []string `json:"parameters"`
	Returns    []string `json:"returns"`
	NatSpec    string   `json:"natspec"`
}

type StateVar struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Visibility string `json:"visibility"`
}

type Detector struct {
	Check       string `json:"check"`
	Impact      string `json:"impact"`
	Confidence  string `json:"confidence"`
	Description string `json:"description"`
}

type slitherRaw struct {
	Success bool `json:"success"`
	Results struct {
		Detectors []struct {
			Check       string `json:"check"`
			Impact      string `json:"impact"`
			Confidence  string `json:"confidence"`
			Description string `json:"description"`
			Elements    []struct {
				SourceMapping struct {
					Filename string `json:"filename_relative"` // <-- key field
				} `json:"source_mapping"`
			} `json:"elements"`
		} `json:"detectors"`
		Contracts []struct {
			Name          string `json:"name"`
			SourceMapping struct {
				Filename string `json:"filename_relative"` // <-- key field
			} `json:"source_mapping"`
			Inheritance []struct {
				Name string `json:"name"`
			} `json:"inheritance"`
			Functions []struct {
				Name       string `json:"name"`
				Visibility string `json:"visibility"`
				Modifiers  []struct {
					Name string `json:"name"`
				} `json:"modifiers"`
				Parameters []struct {
					Type string `json:"type"`
				} `json:"parameters"`
				ReturnType []struct {
					Type string `json:"type"`
				} `json:"returns"`
			} `json:"functions"`
			StateVariables []struct {
				Name       string `json:"name"`
				Type       string `json:"type"`
				Visibility string `json:"visibility"`
			} `json:"state_variables"`
		} `json:"contracts"`
	} `json:"results"`
}
