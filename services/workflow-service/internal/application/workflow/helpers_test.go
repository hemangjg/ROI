package workflow

import "testing"

func TestOptionalUUID(t *testing.T) {
	if got := optionalUUID(""); got != nil {
		t.Fatalf("empty => nil, got %v", got)
	}
	if got := optionalUUID("not-a-uuid"); got != nil {
		t.Fatalf("invalid => nil, got %v", got)
	}
	const raw = "11111111-1111-1111-1111-111111111111"
	got := optionalUUID(raw)
	if got == nil || got.String() != raw {
		t.Fatalf("valid uuid: got %v", got)
	}
}

func TestParseUUID(t *testing.T) {
	if _, err := parseUUID("bad"); err == nil {
		t.Fatal("expected error for invalid uuid")
	}
	id, err := parseUUID("22222222-2222-2222-2222-222222222222")
	if err != nil || !id.Valid {
		t.Fatalf("parse valid uuid: valid=%v err=%v", id.Valid, err)
	}
}
