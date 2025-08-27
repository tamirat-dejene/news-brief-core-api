package dto

// CreateUserRequest is the DTO for creating a new user.
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32,containsuppercase,containslowercase,containsdigit,containssymbol"`
	Fullname string `json:"fullname" binding:"required,min=3,max=50"`
}

// LoginRequest is the DTO for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest is the DTO for user registration.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
	Fullname string `json:"fullname" binding:"required,min=3,max=50"`
}

// UpdateUserRequest is the DTO for updating user profile.
type UpdateUserRequest struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=32"`
	Fullname  *string `json:"fullname,omitempty" binding:"omitempty,max=50"`
}

// ForgotPasswordRequest is the DTO for requesting password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest is the DTO for resetting password.
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Verifier string `json:"verifier" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

// VerifyEmailRequest is the DTO for verifying email.
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// ResendVerificationRequest is the DTO for resending verification email.
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// RefreshTokenRequest is the DTO for refreshing access tokens.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AddSubscriptionRequest struct {
	SourceKey string `json:"source_key" binding:"required"`
}

type ReplaceTopicsRequest struct {
	Topics []string `json:"topics" binding:"required,dive,min=1"`
}

type UpdateTopicsRequest struct {
	Action string   `json:"action" binding:"required,oneof=add remove"`
	Topics []string `json:"topics" binding:"required,dive,min=1"`
}

// NotificationsRequestDTO defines the nested notifications object for preference updates.
type NotificationsRequestDTO struct {
	DailyBrief   *bool `json:"daily_brief"`
	BreakingNews *bool `json:"breaking_news"`
}

// UpdatePreferencesRequest defines the body for the PATCH /v1/me/preferences endpoint.
type UpdatePreferencesRequest struct {
	Lang          *string                  `json:"lang"`
	BriefType     *string                  `json:"brief_type"`
	DataSaver     *bool                    `json:"data_saver"`
	Notifications *NotificationsRequestDTO `json:"notifications,omitempty"`
}
