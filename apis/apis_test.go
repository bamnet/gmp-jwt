package apis

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestLookup(t *testing.T) {
	for _, row := range []struct {
		apiName      []string
		wantScope    string
		wantAudience string
	}{
		{[]string{"routes"}, RoutesScope, RoutesAudience},
		{[]string{"addressvalidation"}, AddressValidationScope, AddressValidationAudience},
		{[]string{"addressvalidation", "routes"}, fmt.Sprintf("%s %s", AddressValidationScope, RoutesScope), ""},
		{[]string{"unknown"}, "", ""},
	} {
		got := Lookup(row.apiName)
		if got.Audience != row.wantAudience {
			t.Errorf("Lookup(%s).Audience got %s, want %s", row.apiName, got.Audience, row.wantAudience)
		}
		if got.Scope != row.wantScope {
			t.Errorf("Lookup(%s).Scope got %s, want %s", row.apiName, got.Scope, row.wantScope)
		}
	}
}

func TestWildcardLookup(t *testing.T) {
	got := Lookup([]string{"*"})

	wantAudience := ""
	if got.Audience != wantAudience {
		t.Errorf("Lookup(*).Audience got %s, want %s", got.Audience, wantAudience)
	}

	wantScope := []string{RoutesScope, AddressValidationScope}
	gotScope := strings.Split(got.Scope, " ")

	if diff := cmp.Diff(wantScope, gotScope, cmpopts.SortSlices(func(x, y string) bool { return x > y })); diff != "" {
		t.Errorf("Lookup(*).Scope mismatch (-want +got):\n%s", diff)
	}
}
