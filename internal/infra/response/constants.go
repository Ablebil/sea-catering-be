package response

// Auth Domain
const (
	EmailAlreadyExists  = "Email already exists"
	UserNotFound        = "User not found"
	InvalidOTP          = "Invalid or expired OTP"
	InvalidCredentials  = "Invalid email or password"
	InvalidRefreshToken = "Invalid or expired refresh token"
	UserNotVerified     = "User not verified"

	FailedFindUser           = "Failed to find user"
	FailedCreateUser         = "Failed to create user"
	FailedUpdateUser         = "Failed to update user"
	FailedAddRefreshToken    = "Failed to add refresh token"
	FailedGetRefreshTokens   = "Failed to get refresh tokens"
	FailedRemoveRefreshToken = "Failed to remove refresh token"

	RegisterSuccess     = "Registration successful. OTP has been sent to email"
	VerifyOTPSuccess    = "Verification successful"
	LoginSuccess        = "Login successful"
	RefreshTokenSuccess = "Token refresh successful"
	LogoutSuccess       = "Logout successful"
)

// Others
const (
	FailedHashPassword         = "Failed to hash password"
	FailedGenerateOTP          = "Failed to generate OTP"
	FailedStoreOTP             = "Failed to store OTP"
	FailedDeleteOTP            = "Failed to delete OTP"
	FailedSendOTPEmail         = "Failed to send OTP email"
	FailedGenerateRefreshToken = "Failed to generate refresh token"
	FailedGenerateAccessToken  = "Failed to generate access token"
)

// Handler
const (
	FailedParsingRequestBody = "Failed parsing request body"
	FailedValidateRequest    = "Failed to validate request"
)
