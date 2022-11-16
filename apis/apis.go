// Package apis contains JWT claim values for Google Maps Platform APIs.
package apis

import (
	"strings"
)

const RoutesScope = "https://www.googleapis.com/auth/geo-platform.routes"
const RoutesAudience = "https://routes.googleapis.com/"

const AddressValidationScope = "https://www.googleapis.com/auth/maps-platform.addressvalidation"
const AddressValidationAudience = "https://addressvalidation.googleapis.com/"

// APITokenInfo represents the claims a JWT should include to authenticate.
type APITokenInfo struct {
	Scope    string
	Audience string
}

// APIs maps shortnames of each API to information about their scope and audience.
var APIs = map[string]APITokenInfo{
	"routes":            {RoutesScope, RoutesAudience},
	"addressvalidation": {AddressValidationScope, AddressValidationAudience},
}

// Lookup returns scope and audience information suitable for the supplied APIs.
func Lookup(names []string) APITokenInfo {
	scopes := []string{}
	audience := []string{}
	if names[0] == "*" {
		for _, api := range APIs {
			scopes = append(scopes, api.Scope)
			audience = append(audience, api.Audience)
		}
	} else {
		for _, name := range names {
			if api, ok := APIs[name]; ok {
				scopes = append(scopes, api.Scope)
				audience = append(audience, api.Audience)
			}
		}
	}

	token := APITokenInfo{}
	// An audience constraint can only be applied when the token has 1 audience.
	if len(audience) == 1 {
		token.Audience = audience[0]
	}
	// Join together scopes with a space.
	token.Scope = strings.Join(scopes, " ")

	return token
}
