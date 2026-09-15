package opencode

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

// OpenCode Zen's edge (Cloudflare) rejects chat-completion requests that don't look
// like they came from the OpenCode CLI: a generic User-Agent from a datacenter IP —
// exactly what every VPS running Bifrost presents — gets rejected as a
// FreeUsageLimitError, and some free-tier models 400 outright without an
// x-opencode-session header. None of this is documented by OpenCode; it was
// reverse-engineered by github.com/diegosouzapw/OmniRoute, which spoofs the CLI's
// own request shape to get through. This mirrors that fix for Bifrost's Zen provider.
const (
	defaultCLIUserAgent = "opencode"
	defaultCLIClient    = "desktop"
	defaultCLIProject   = "global"
)

// cliHeadersEnabled reports whether OpenCode CLI identity headers should be
// synthesized. Opt out with OPENCODE_SYNTHESIZE_CLI_HEADERS=false if a future
// upstream change makes this unnecessary or counterproductive.
func cliHeadersEnabled() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("OPENCODE_SYNTHESIZE_CLI_HEADERS")))
	return v != "0" && v != "false" && v != "no" && v != "off"
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

// buildCLIHeaders returns extraHeaders merged with synthesized OpenCode CLI identity
// headers. Caller-configured values in extraHeaders always win (merged in last), so a
// real "opencode-cli/x.y.z" User-Agent or a specific client/project id can still be
// set explicitly via the provider's network_config.extra_headers. x-opencode-request
// and x-opencode-session are always freshly generated — reusing them across requests
// is exactly what would make Bifrost's traffic look scripted instead of CLI-like.
//
// ponytail: x-opencode-session is a fresh UUID per request rather than a
// conversation-stable fingerprint (OmniRoute hashes model+system+first-user-message so
// repeated turns of the same conversation hit the upstream's prompt cache). Add that if
// this gateway ever serves a workload where OpenCode Zen's prompt-caching discount
// matters — for a single self-hosted agent it doesn't.
func buildCLIHeaders(extraHeaders map[string]string) map[string]string {
	if !cliHeadersEnabled() {
		return extraHeaders
	}

	headers := map[string]string{
		"User-Agent":         envOr("OPENCODE_USER_AGENT", defaultCLIUserAgent),
		"x-opencode-client":  envOr("OPENCODE_CLIENT", defaultCLIClient),
		"x-opencode-project": envOr("OPENCODE_PROJECT", defaultCLIProject),
		"x-opencode-request": uuid.NewString(),
		"x-opencode-session": uuid.NewString(),
	}
	for k, v := range extraHeaders {
		headers[k] = v
	}
	return headers
}
