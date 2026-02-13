package azure

// SignInActivity represents Microsoft Graph's `signInActivity` object returned on the user entity.
type SignInActivity struct {
	LastSignInDateTime                *string `json:"lastSignInDateTime,omitempty"`
	LastSignInRequestId               *string `json:"lastSignInRequestId,omitempty"`
	LastNonInteractiveSignInDateTime  *string `json:"lastNonInteractiveSignInDateTime,omitempty"`
	LastNonInteractiveSignInRequestId *string `json:"lastNonInteractiveSignInRequestId,omitempty"`
	LastSuccessfulSignInDateTime      *string `json:"lastSuccessfulSignInDateTime,omitempty"`
	LastSuccessfulSignInRequestId     *string `json:"lastSuccessfulSignInRequestId,omitempty"`
}
