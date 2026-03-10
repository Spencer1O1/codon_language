package validator

// PatchDoc represents the mutation.patch.patch_definition structure.
type PatchDoc struct {
	Summary     string      `json:"summary" yaml:"summary"`
	Rationale   string      `json:"rationale" yaml:"rationale"`
	Risk        string      `json:"risk" yaml:"risk"`
	Evidence    interface{} `json:"evidence" yaml:"evidence"`
	Warnings    interface{} `json:"warnings" yaml:"warnings"`
	Suggestions interface{} `json:"suggestions" yaml:"suggestions"`
	Confidence  string      `json:"confidence" yaml:"confidence"`
	Changes     []Change    `json:"changes" yaml:"changes"`
}

type Change struct {
	ID         string      `json:"id" yaml:"id"`
	Target     string      `json:"target" yaml:"target"`
	Operations []Operation `json:"operations" yaml:"operations"`
}

type Operation struct {
	ID       string      `json:"id" yaml:"id"`
	Op       string      `json:"op" yaml:"op"`
	Path     string      `json:"path" yaml:"path"`
	Value    interface{} `json:"value" yaml:"value"`
	OldValue interface{} `json:"old_value" yaml:"old_value"`
	Required *bool       `json:"required" yaml:"required"`
	Reason   string      `json:"reason" yaml:"reason"`
}
