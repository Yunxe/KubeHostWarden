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
