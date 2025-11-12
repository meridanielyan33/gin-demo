package errors

const (
	AuthHeaderMissing    = "Authorization header is missing"
	LoggedOut            = "You have been logged out"
	InvalidToken         = "Invalid or missing token"
	ExpiredToken         = "Token is expired"
	RedisTokenFail       = "Failed to check user's token in Redis"
	TokenDeleteFailRedis = "Failed to delete token from Redis"
	InvalidReqData       = "Invalid request data"
	UserNotAuthenticated = "User not authenticated"
)
