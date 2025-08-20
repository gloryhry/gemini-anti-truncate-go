package gemini

// GenerateContentRequest represents the request body for the generateContent endpoint.
type GenerateContentRequest struct {
	Contents           []Content          `json:"contents"`
	SystemInstruction  *SystemInstruction `json:"systemInstruction,omitempty"`  // Kept for compatibility
	SystemInstruction_ *SystemInstruction `json:"system_instruction,omitempty"` // Official field
	GenerationConfig   *GenerationConfig  `json:"generationConfig,omitempty"`
	SafetySettings     []SafetySetting    `json:"safetySettings,omitempty"`
	Tools              []Tool             `json:"tools,omitempty"`
}

// GetSystemInstruction provides a unified way to get the system instruction,
// checking both snake_case and camelCase fields for maximum compatibility.
func (r *GenerateContentRequest) GetSystemInstruction() *SystemInstruction {
	if r.SystemInstruction_ != nil {
		return r.SystemInstruction_
	}
	return r.SystemInstruction
}

// SetSystemInstruction provides a unified way to set the system instruction.
// It clears the camelCase field to avoid sending both.
func (r *GenerateContentRequest) SetSystemInstruction(si *SystemInstruction) {
	r.SystemInstruction_ = si
	r.SystemInstruction = nil
}

// SystemInstruction represents the content of a system instruction.
type SystemInstruction struct {
	Role  string `json:"role,omitempty"`
	Parts []Part `json:"parts"`
}

// GenerateContentResponse represents the full response for a non-streaming request,
// or a single chunk in a streaming request.
type GenerateContentResponse struct {
	Candidates     []Candidate    `json:"candidates"`
	PromptFeedback PromptFeedback `json:"promptFeedback,omitempty"`
}

// Candidate represents a single response candidate from the model.
type Candidate struct {
	Content          Content           `json:"content"`
	FinishReason     string            `json:"finishReason,omitempty"`
	Index            int               `json:"index"`
	SafetyRatings    []SafetyRating    `json:"safetyRatings,omitempty"`
	CitationMetadata *CitationMetadata `json:"citationMetadata,omitempty"`
}

// Content represents a message in the conversation history.
type Content struct {
	Parts []Part `json:"parts"`
	Role  string `json:"role"` // "user", "model", or "tool"
}

// Part represents a single part of a Content message.
type Part struct {
	Text         string        `json:"text,omitempty"`
	FunctionCall *FunctionCall `json:"functionCall,omitempty"`
	// This field is used internally by the proxy to identify and handle "thought" blocks from the model.
	Thought bool `json:"thought,omitempty"`
}

// FunctionCall represents a function call requested by the model.
type FunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

// GenerationConfig specifies model parameters for the generation.
type GenerationConfig struct {
	Temperature      float64     `json:"temperature,omitempty"`
	TopP             float64     `json:"topP,omitempty"`
	TopK             int         `json:"topK,omitempty"`
	MaxOutputTokens  int         `json:"maxOutputTokens,omitempty"`
	CandidateCount   int         `json:"candidateCount,omitempty"`
	StopSequences    []string    `json:"stopSequences,omitempty"`
	ResponseMIMEType string      `json:"responseMimeType,omitempty"`
	ResponseSchema   interface{} `json:"responseSchema,omitempty"` // Can be complex, so interface{} is used.
}

// SafetySetting configures the safety thresholds for different categories.
type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// Tool represents a tool definition that the model can use.
type Tool struct {
	FunctionDeclarations []interface{} `json:"functionDeclarations,omitempty"` // Can be complex.
}

// PromptFeedback provides feedback on the prompt, such as safety ratings.
type PromptFeedback struct {
	BlockReason   string         `json:"blockReason,omitempty"`
	SafetyRatings []SafetyRating `json:"safetyRatings,omitempty"`
}

// SafetyRating provides the safety rating for a piece of content.
type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

// CitationMetadata contains citation information for the generated content.
type CitationMetadata struct {
	CitationSources []CitationSource `json:"citationSources"`
}

// CitationSource provides a single citation source with its URI and license.
type CitationSource struct {
	StartIndex int    `json:"startIndex,omitempty"`
	EndIndex   int    `json:"endIndex,omitempty"`
	URI        string `json:"uri,omitempty"`
	License    string `json:"license,omitempty"`
}

// ErrorResponse represents a standard error from the Gemini API.
type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}
