package apis

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var sort = cmpopts.SortSlices(func(x, y string) bool { return x > y })

func TestLookup(t *testing.T) {
	for _, row := range []struct {
		apiName      []string
		wantScope    []string
		wantAudience string
		wantError    bool
	}{
		{[]string{"routes"}, []string{RoutesScope}, RoutesAudience, false},
		{[]string{"addressvalidation"}, []string{AddressValidationScope}, AddressValidationAudience, false},
		{[]string{"addressvalidation", "routes"}, []string{AddressValidationScope, RoutesScope}, "", false},
		{[]string{"airquality"}, []string{GCPScope}, AirQualityAudience, false},
		{[]string{"airquality", "solar"}, []string{""}, "", true},
		{[]string{"solar", "places"}, []string{""}, "", true},
		{[]string{"unknown"}, []string{""}, "", false},
		{[]string{}, []string{""}, "", false},
	} {
		gotToken, gotError := Lookup(row.apiName)
		if row.wantError {
			if gotError == nil {
				t.Errorf("Lookup(%s) got error %v, wanted %v", row.apiName, gotError, row.wantError)
			}
			continue
		}
		if gotToken.Audience != row.wantAudience {
			t.Errorf("Lookup(%s).Audience got %s, want %s", row.apiName, gotToken.Audience, row.wantAudience)
		}

		if diff := cmp.Diff(row.wantScope, strings.Split(gotToken.Scope, " "), sort); diff != "" {
			t.Errorf("Lookup(%s).Scope mismatch (-want +got):\n%s", row.apiName, diff)
		}
	}
}

func TestWildcardLookup(t *testing.T) {
	got, _ := Lookup([]string{"*"})

	wantAudience := ""
	if got.Audience != wantAudience {
		t.Errorf("Lookup(*).Audience got %s, want %s", got.Audience, wantAudience)
	}

	wantScope := []string{RoutesScope, AddressValidationScope, PlacesScope}
	gotScope := strings.Split(got.Scope, " ")

	if diff := cmp.Diff(wantScope, gotScope, sort); diff != "" {
		t.Errorf("Lookup(*).Scope mismatch (-want +got):\n%s", diff)
	}
}
