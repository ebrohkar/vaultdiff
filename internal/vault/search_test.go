package vault

import (
	"testing"
)

func TestSearchOptions_DefaultMount(t *testing.T) {
	opts := SearchOptions{Keyword: "db"}
	if opts.Mount != "" {
		t.Errorf("expected empty mount before normalisation, got %q", opts.Mount)
	}
}

func TestSearchResult_Fields(t *testing.T) {
	r := SearchResult{
		Path:        "secret/data/myapp/config",
		MatchedKeys: []string{"db_password", "db_host"},
	}
	if r.Path != "secret/data/myapp/config" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if len(r.MatchedKeys) != 2 {
		t.Errorf("expected 2 matched keys, got %d", len(r.MatchedKeys))
	}
}

func TestSearchSecrets_EmptyKeyword(t *testing.T) {
	client := &Client{}
	_, err := SearchSecrets(client, SearchOptions{Mount: "secret", Keyword: ""})
	if err == nil {
		t.Fatal("expected error for empty keyword")
	}
	expected := "search keyword must not be empty"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestMatchingKeys_CaseSensitive(t *testing.T) {
	data := map[string]interface{}{
		"DB_PASSWORD": "s3cr3t",
		"api_key":     "abc123",
		"host":        "localhost",
	}
	matched := matchingKeys(data, "DB", true)
	if len(matched) != 1 || matched[0] != "DB_PASSWORD" {
		t.Errorf("expected [DB_PASSWORD], got %v", matched)
	}
}

func TestMatchingKeys_CaseInsensitive(t *testing.T) {
	data := map[string]interface{}{
		"DB_PASSWORD": "s3cr3t",
		"api_key":     "abc123",
		"host":        "localhost",
	}
	matched := matchingKeys(data, "db", false)
	if len(matched) != 1 {
		t.Errorf("expected 1 match, got %d: %v", len(matched), matched)
	}
}

func TestMatchingKeys_ValueMatch(t *testing.T) {
	data := map[string]interface{}{
		"username": "admin",
		"token":    "supersecret-token",
	}
	matched := matchingKeys(data, "supersecret", true)
	if len(matched) != 1 || matched[0] != "token" {
		t.Errorf("expected [token], got %v", matched)
	}
}

func TestMatchingKeys_NoMatch(t *testing.T) {
	data := map[string]interface{}{
		"host": "localhost",
		"port": "5432",
	}
	matched := matchingKeys(data, "aws", false)
	if len(matched) != 0 {
		t.Errorf("expected no matches, got %v", matched)
	}
}

func TestMatchingKeys_EmptyData(t *testing.T) {
	matched := matchingKeys(map[string]interface{}{}, "keyword", false)
	if len(matched) != 0 {
		t.Errorf("expected empty result, got %v", matched)
	}
}
