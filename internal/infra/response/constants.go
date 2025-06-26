package response

// Auth Domain
const (
	EmailAlreadyExists  = "Email already exists"
	UserNotFound        = "User not found"
	InvalidOTP          = "Invalid or expired OTP"
	InvalidCredentials  = "Invalid email or password"
	InvalidRefreshToken = "Invalid or expired refresh token"
	UserNotVerified     = "User not verified"
	OAuthStateNotFound  = "OAuth state not found"
	OAuthStateInvalid   = "OAuth state invalid"

	FailedFindUser           = "Failed to find user"
	FailedCreateUser         = "Failed to create user"
	FailedUpdateUser         = "Failed to update user"
	FailedAddRefreshToken    = "Failed to add refresh token"
	FailedGetRefreshTokens   = "Failed to get refresh tokens"
	FailedRemoveRefreshToken = "Failed to remove refresh token"
	FailedExchangeOAuthToken = "Failed to exchange OAuth token"
	FailedGetOAuthProfile    = "Failed to get OAuth profile"
	FailedGoogleLogin        = "Failed to initiate Google login"

	RegisterSuccess     = "Registration successful. OTP has been sent to email"
	VerifyOTPSuccess    = "Verification successful"
	LoginSuccess        = "Login successful"
	RefreshTokenSuccess = "Token refresh successful"
	LogoutSuccess       = "Logout successful"
)

// Testimonial Domain
const (
	GetAllTestimonialsSuccess = "Get all testimonials successful"
	CreateTestimonialSuccess  = "Create testimonial successful"

	FailedGetAllTestimonials = "Failed to get all testimonials"
	FailedCreateTestimonial  = "Failed to create testimonial"
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
	FailedGenerateOAuthState   = "Failed to generate OAuth state"
	FailedStoreOAuthState      = "Failed to store OAuth state"
	FailedDeleteOAuthState     = "Failed to delete OAuth state"
	FailedGenerateOAuthLink    = "Failed to generate OAuth link"
	FailedOAuthCallback        = "Failed to handle OAuth callback"
	FailedUploadFile           = "Failed to upload file"
)

// Handler
const (
	FailedParsingRequestBody    = "Failed parsing request body"
	FailedValidateRequest       = "Failed to validate request"
	MissingAccessToken          = "Missing access token"
	InvalidAccessToken          = "Invalid access token"
	InvalidOrMissingBearerToken = "Invalid or missing bearer token"
	InvalidFormData             = "Invalid form data"
	FileIsRequired              = "File is required"
	FailedToOpenFile            = "Failed to open file"
)
