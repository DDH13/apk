package authentication

import (
	"encoding/json"

	"github.com/wso2/apk/gateway/enforcer/internal/dto"
	"github.com/wso2/apk/gateway/enforcer/internal/requestconfig"
	"github.com/wso2/apk/gateway/enforcer/internal/transformer"
)

// ValidateToken validates the JWT token.
func ValidateToken(rch *requestconfig.Holder, jwtTransformer *transformer.JWTTransformer) *dto.ImmediateResponse {
	if rch != nil {
		if rch.ExternalProcessingEnvoyMetadata.JwtAuthenticationData != nil {
			jwtValidationInfo := jwtTransformer.TransformJWTClaims(rch.MatchedAPI.OrganizationID, rch.ExternalProcessingEnvoyMetadata.JwtAuthenticationData)
			if jwtValidationInfo != nil {
				if jwtValidationInfo.Valid {
					rch.JWTValidationInfo = jwtValidationInfo
					return nil
				}
			}
		}
	}
	errorResponse := &dto.ErrorResponse{ErrorMessage: "Invalid Credentials", Code: "900901", ErrorDescription: "Make sure you have provided the correct security credentials"}
	jsonData, _ := json.MarshalIndent(errorResponse, "", "  ")
	return &dto.ImmediateResponse{StatusCode: 401, Message: string(jsonData)}
}
