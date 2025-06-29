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

// User Domain
const (
	FailedRemoveUnverifiedUsers = "Failed to remove unverified users"
	FailedGetUserProfile        = "Failed to get user profile"

	GetProfileSuccess = "Get profile successful"
)

// Testimonial Domain
const (
	FailedGetAllTestimonials = "Failed to get all testimonials"
	FailedCreateTestimonial  = "Failed to create testimonial"

	GetAllTestimonialsSuccess = "Get all testimonials successful"
	CreateTestimonialSuccess  = "Create testimonial successful"
)

// Meal Plan Domain
const (
	MealPlanNotFound = "Meal plan not found"

	FailedGetAllMealPlans = "Failed to get all meal plans"
	FailedGetMealPlanByID = "Failed to get meal plan by ID"
	FailedCreateMealPlan  = "Failed to create meal plan"

	GetAllMealPlansSuccess = "Get all meal plans successful"
	GetMealPlanByIDSuccess = "Get meal plan by ID successful"
	CreateMealPlanSuccess  = "Create meal plan successful"
)

// Subscription Domain
const (
	SubscriptionNotFound = "Subscription not found"

	FailedSaveSubscription            = "Failed to save subscription"
	FailedCreatePaymentTransaction    = "Failed to create payment transaction"
	FailedGetAllSubscriptions         = "Failed to get all subscriptions"
	FailedGetSubscriptionByID         = "Failed to get subscription by ID"
	FailedUpdateSubscription          = "Failed to update subscription"
	FailedPauseSubscription           = "Failed to pause subscription"
	FailedCancelSubscription          = "Failed to cancel subscription"
	FailedGetExpiredSubscriptions     = "Failed to get expired subscriptions"
	FailedGetNewSubscriptionsCount    = "Failed to get new subscriptions count"
	FailedCalculateMMR                = "Failed to calculate MMR"
	FailedGetTotalActiveSubscriptions = "Failed to get total active subscriptions"
	FailedGetReactivationStats        = "Failed to get reactivation stats"

	CreateSubscriptionSuccess          = "Subscription created successful"
	GetAllSubscriptionsSuccess         = "Get all subscriptions successful"
	PauseSubscriptionSuccess           = "Subscription paused successful"
	CancelSubscriptionSuccess          = "Subscription cancelled successful"
	GetNewSubscriptionsStatsSuccess    = "Get new subscriptions stats success"
	GetMRRStatsSuccess                 = "Get MRR stats success"
	GetTotalActiveSubscriptionsSuccess = "Get total active subscriptions success"
	GetReactivationStatsSuccess        = "Get reactivation stats success"
	WebhookProcessedSuccess            = "Webhook processed successful"
)

// Others
const (
	FailedHashPassword           = "Failed to hash password"
	FailedGenerateOTP            = "Failed to generate OTP"
	FailedStoreOTP               = "Failed to store OTP"
	FailedDeleteOTP              = "Failed to delete OTP"
	FailedSendOTPEmail           = "Failed to send OTP email"
	FailedGenerateRefreshToken   = "Failed to generate refresh token"
	FailedGenerateAccessToken    = "Failed to generate access token"
	FailedGenerateOAuthState     = "Failed to generate OAuth state"
	FailedStoreOAuthState        = "Failed to store OAuth state"
	FailedDeleteOAuthState       = "Failed to delete OAuth state"
	FailedGenerateOAuthLink      = "Failed to generate OAuth link"
	FailedOAuthCallback          = "Failed to handle OAuth callback"
	FailedUploadFile             = "Failed to upload file"
	FailedReadFileForValidation  = "Failed read file for validation"
	FailedToResetFileReadPointer = "Failed to reset file read pointer"
	FileSizeExceedsLimit         = "File size exceeds the limit"
	InvalidFileType              = "Invalid file type. Only JPG, JPEG, and PNG are allowed."
	InvalidOrderID               = "Invalid order ID"
	InvalidTransactionStatus     = "Invalid transaction status"
)

// Handler
const (
	FailedParsingRequestBody    = "Failed parsing request body"
	FailedParsingRequestParams  = "Failed parsing request params"
	FailedValidateRequest       = "Failed to validate request"
	MissingAccessToken          = "Missing access token"
	InvalidAccessToken          = "Invalid access token"
	InvalidOrMissingBearerToken = "Invalid or missing bearer token"
	InvalidFormData             = "Invalid form data"
	FileIsRequired              = "File is required"
	FailedToOpenFile            = "Failed to open file"
	InvalidMealPlanID           = "Invalid meal plan ID"
	InvalidSubscriptionID       = "Invalid subscription ID"
	AdminAccessRequired         = "Admin access required"
)
