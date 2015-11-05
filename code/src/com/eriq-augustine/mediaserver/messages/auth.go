package messages;

import (
   "com/eriq-augustine/mediaserver/util/errors"
);

type RejectedToken struct {
   Success bool
   ReasonCode int
   ReasonDescription string
}

func NewRejectedToken(err errors.TokenValidationError) *RejectedToken {
   return &RejectedToken{false, err.Reason, err.Description()};
}

type AuthorizedToken struct {
   Success bool
   Token string
}

func NewAuthorizedToken(token string) *AuthorizedToken {
   return &AuthorizedToken{true, token};
}
