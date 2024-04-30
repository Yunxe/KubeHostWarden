package constant

type ContextKey string

const (
	UserIDKey    ContextKey = "userId"
	UserEmailKey ContextKey = "userEmail"
)

const (
	DARWIN = "darwin"
	LINUX  = "linux"
)

type ThresholdType string

const (
	ABOVE ThresholdType = "above"
	BELOW ThresholdType = "below"
)
