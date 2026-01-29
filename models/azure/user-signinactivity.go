package azure

import "encoding/json"

type userAlias User

type userUnmarshalJSON struct {
	*userAlias
	SignInActivity *SignInActivity `json:"signInActivity,omitempty"`
}

func (s *User) UnmarshalJSON(data []byte) error {
	aux := userUnmarshalJSON{userAlias: (*userAlias)(s)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if s.LastSuccessfulSignInDateTime == "" && aux.SignInActivity != nil && aux.SignInActivity.LastSuccessfulSignInDateTime != nil {
		s.LastSuccessfulSignInDateTime = *aux.SignInActivity.LastSuccessfulSignInDateTime
	}

	return nil
}

// SignInActivity represents Microsoft Graph's `signInActivity` object returned on the user entity.
type SignInActivity struct {
	LastSignInDateTime                *string `json:"lastSignInDateTime,omitempty"`
	LastSignInRequestId               *string `json:"lastSignInRequestId,omitempty"`
	LastNonInteractiveSignInDateTime  *string `json:"lastNonInteractiveSignInDateTime,omitempty"`
	LastNonInteractiveSignInRequestId *string `json:"lastNonInteractiveSignInRequestId,omitempty"`
	LastSuccessfulSignInDateTime      *string `json:"lastSuccessfulSignInDateTime,omitempty"`
	LastSuccessfulSignInRequestId     *string `json:"lastSuccessfulSignInRequestId,omitempty"`
}
