// Package apis contains JWT claim values for Google Maps Platform APIs.
package apis

import (
	"errors"
	"strings"
)

const RoutesScope = "https://www.googleapis.com/auth/geo-platform.routes"
const RoutesAudience = "https://routes.googleapis.com/"

const AddressValidationScope = "https://www.googleapis.com/auth/maps-platform.addressvalidation"
const AddressValidationAudience = "https://addressvalidation.googleapis.com/"

const PlacesScope = "https://www.googleapis.com/auth/maps-platform.places"
const PlacesAudience = "https://places.googleapis.com/"

const GCPScope = "https://www.googleapis.com/auth/cloud-platform"
const AirQualityAudience = "https://airquality.googleapis.com/"
const SolarAudience = "https://solar.googleapis.com/"

// APITokenInfo represents the claims a JWT should include to authenticate.
type APITokenInfo struct {
	Scope    string
	Audience string
}

// APIs maps shortnames of each API to information about their scope and audience.
var APIs = map[string]APITokenInfo{
	"routes":            {RoutesScope, RoutesAudience},
	"addressvalidation": {AddressValidationScope, AddressValidationAudience},
	"places":            {PlacesScope, PlacesAudience},
	"airquality":        {GCPScope, AirQualityAudience},
	"solar":             {GCPScope, SolarAudience},
}

var ErrIncompatibleAPIs = errors.New("request cannot include a multiple gcp-scoped apis or a mix of gcp-scoped and maps-scoped")

// Lookup returns scope and audience information suitable for the supplied APIs.
func Lookup(names []string) (*APITokenInfo, error) {
	scopes := map[string]bool{}
	audience := []string{}
	token := APITokenInfo{}

	if len(names) == 0 {
		return &token, nil
	}

	if IsWildcard(names) {
		for _, api := range APIs {
			if api.Scope == GCPScope {
				// Skip the GCPScope, it should be part of a wildcard.
				continue
			}
			scopes[api.Scope] = true
			audience = append(audience, api.Audience)
		}
	} else {
		for _, name := range names {
			if api, ok := APIs[name]; ok {
				scopes[api.Scope] = true
				audience = append(audience, api.Audience)
			}
		}
	}

	// An audience constraint can only be applied when the token has 1 audience.
	if len(audience) == 1 {
		token.Audience = audience[0]
	}

	// Never return a token that has the GCP scope and no audience.
	// That token would grant too much access.
	if token.Audience == "" {
		if _, found := scopes[GCPScope]; found {
			return nil, ErrIncompatibleAPIs
		}
	}

	// Join together scopes with a space.
	s := []string{}
	for scope := range scopes {
		s = append(s, scope)
	}
	token.Scope = strings.Join(s, " ")

	return &token, nil
}

// IsWildcard identifies if list of API names actually means all APIs.
func IsWildcard(apis []string) bool {
	return len(apis) >= 1 && apis[0] == "*"
}
