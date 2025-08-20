package gemini

const (
	FinishToken          = "[RESPONSE_FINISHED]"
	UserPromptSuffix     = "\n\n(Note: If you are done, please end your response with " + FinishToken + ")"
	RetryPrompt          = "Please continue generating the response from where you left off. Do not repeat the previous content."
	DefaultUpstreamURL   = "https://generativelanguage.googleapis.com"
	DefaultMaxRetries    = 20
	DefaultHTTPPort      = 8080
	TokenLookbehindChars = len(FinishToken) + 5 // A little buffer for lookbehind
)

var TargetModels = []string{
	"gemini-1.5-pro-latest",
	"gemini-1.5-flash-latest",
	"gemini-pro",
}

var RetryableStatus = []int{503, 403, 429}
var FatalStatus = []int{500}
