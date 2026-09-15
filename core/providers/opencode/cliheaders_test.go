package opencode

import "testing"

func TestBuildCLIHeadersDefaults(t *testing.T) {
	headers := buildCLIHeaders(nil)

	if headers["User-Agent"] != defaultCLIUserAgent {
		t.Errorf("expected User-Agent %q, got %q", defaultCLIUserAgent, headers["User-Agent"])
	}
	if headers["x-opencode-client"] != defaultCLIClient {
		t.Errorf("expected x-opencode-client %q, got %q", defaultCLIClient, headers["x-opencode-client"])
	}
	if headers["x-opencode-project"] != defaultCLIProject {
		t.Errorf("expected x-opencode-project %q, got %q", defaultCLIProject, headers["x-opencode-project"])
	}
	if headers["x-opencode-request"] == "" {
		t.Error("expected x-opencode-request to be set")
	}
	if headers["x-opencode-session"] == "" {
		t.Error("expected x-opencode-session to be set")
	}
}

func TestBuildCLIHeadersFreshPerCall(t *testing.T) {
	first := buildCLIHeaders(nil)
	second := buildCLIHeaders(nil)

	if first["x-opencode-request"] == second["x-opencode-request"] {
		t.Error("expected x-opencode-request to differ across calls")
	}
	if first["x-opencode-session"] == second["x-opencode-session"] {
		t.Error("expected x-opencode-session to differ across calls")
	}
}

func TestBuildCLIHeadersCallerOverridesWin(t *testing.T) {
	headers := buildCLIHeaders(map[string]string{
		"User-Agent": "opencode-cli/1.2.3",
		"X-Custom":   "kept",
	})

	if headers["User-Agent"] != "opencode-cli/1.2.3" {
		t.Errorf("expected caller-configured User-Agent to win, got %q", headers["User-Agent"])
	}
	if headers["X-Custom"] != "kept" {
		t.Error("expected unrelated caller header to be preserved")
	}
}

func TestBuildCLIHeadersDisabledViaEnv(t *testing.T) {
	t.Setenv("OPENCODE_SYNTHESIZE_CLI_HEADERS", "false")

	extra := map[string]string{"X-Only": "this"}
	headers := buildCLIHeaders(extra)

	if len(headers) != 1 || headers["X-Only"] != "this" {
		t.Errorf("expected synthesis disabled to return extraHeaders unchanged, got %v", headers)
	}
}

func TestBuildCLIHeadersEnvOverrides(t *testing.T) {
	t.Setenv("OPENCODE_USER_AGENT", "opencode-cli/9.9.9")
	t.Setenv("OPENCODE_CLIENT", "custom-client")
	t.Setenv("OPENCODE_PROJECT", "custom-project")

	headers := buildCLIHeaders(nil)

	if headers["User-Agent"] != "opencode-cli/9.9.9" {
		t.Errorf("expected env-overridden User-Agent, got %q", headers["User-Agent"])
	}
	if headers["x-opencode-client"] != "custom-client" {
		t.Errorf("expected env-overridden client, got %q", headers["x-opencode-client"])
	}
	if headers["x-opencode-project"] != "custom-project" {
		t.Errorf("expected env-overridden project, got %q", headers["x-opencode-project"])
	}
}
