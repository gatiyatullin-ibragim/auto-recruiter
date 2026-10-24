package contacts

import "testing"

func TestExtractOnlyFoundUsefulEmails(t *testing.T) {
	got := Extract("HR [at] acme [dot] kz; noreply@example.com; evil@other.kz", "https://acme.kz/contacts")
	if len(got) != 2 || got[0].Email != "hr@acme.kz" || got[0].SourceURL == "" {
		t.Fatalf("unexpected contacts: %+v", got)
	}
}
