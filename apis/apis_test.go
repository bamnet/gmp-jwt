package apis

import "testing"

func TestLookup(t *testing.T) {
	for _, row := range []struct {
		apiName      []string
		wantScope    string
		wantAudience string
	}{
		{[]string{"routes"}, RoutesScope, RoutesAudience},
		{[]string{"unknown"}, "", ""},
		{[]string{"*"}, RoutesScope, RoutesAudience},
	} {
		got := Lookup(row.apiName)
		if got.Audience != row.wantAudience {
			t.Errorf("Lookup(%s).Audience got %s, want %s", row.apiName, got.Audience, row.wantAudience)
		}
		if got.Scope != row.wantScope {
			t.Errorf("Lookup(%s).Scope got %s, want %s", row.apiName, got.Audience, row.wantAudience)
		}
	}
}
