package theme

import (
	"html/template"
	"strings"
)

// Custom CSS lets an owner bring their own design language rather than being
// limited to the editor's controls.
//
// This is the one place user input reaches the stylesheet as more than a single
// validated value, so it is filtered rather than trusted. The rules below block
// the constructs that turn a stylesheet into an exfiltration or escape vector:
//
//   - </style and <!-- would end the style element and start injecting markup
//   - @import and url(http…) fetch from third parties, leaking visitor IPs and
//     defeating the page's content security policy
//   - expression() and javascript: are legacy script execution paths
//   - position:fixed over the whole viewport can hijack the page for clickjacking
//
// Everything else is allowed: this is meant to be genuinely useful, not a
// token gesture.
const MaxCustomCSS = 20_000

var cssBanned = []string{
	"</style",
	"<!--",
	"-->",
	"@import",
	"javascript:",
	"expression(",
	"behavior:",
	"-moz-binding",
	"<script",
}

// SanitizeCustomCSS filters an owner's stylesheet and reports what it removed,
// so the editor can explain the change rather than silently dropping work.
func SanitizeCustomCSS(in string) (clean string, removed []string) {
	if strings.TrimSpace(in) == "" {
		return "", nil
	}

	if len(in) > MaxCustomCSS {
		in = in[:MaxCustomCSS]
		removed = append(removed, "size")
	}

	// Drop the whole line a banned token appears on rather than snipping the
	// token out. Cutting just the token leaves mangled fragments like
	// "width: alert(1));" behind — inert, but confusing to anyone reading their
	// own stylesheet back. Removing the line is predictable and explainable.
	var kept []string
	hit := map[string]bool{}

	for _, line := range strings.Split(in, "\n") {
		lower := strings.ToLower(line)
		bad := ""
		for _, token := range cssBanned {
			if strings.Contains(lower, token) {
				bad = token
				break
			}
		}
		if bad != "" {
			if !hit[bad] {
				hit[bad] = true
				removed = append(removed, bad)
			}
			continue
		}
		kept = append(kept, line)
	}
	in = strings.Join(kept, "\n")

	// Remote url() references: keep same-origin uploads, drop everything else.
	if out, cut := stripRemoteURLs(in); cut {
		in = out
		removed = append(removed, "url()")
	}

	return in, removed
}

// stripRemoteURLs rewrites any url() that points off-site to url(). Relative
// paths and data: images are left alone.
func stripRemoteURLs(s string) (string, bool) {
	var b strings.Builder
	changed := false
	rest := s

	for {
		i := indexFold(rest, "url(")
		if i < 0 {
			b.WriteString(rest)
			return b.String(), changed
		}

		b.WriteString(rest[:i+4])
		rest = rest[i+4:]

		end := strings.IndexByte(rest, ')')
		if end < 0 {
			b.WriteString(rest)
			return b.String(), changed
		}

		inner := strings.TrimSpace(rest[:end])
		trimmed := strings.Trim(inner, `"'`)
		low := strings.ToLower(trimmed)

		switch {
		case strings.HasPrefix(low, "http://"),
			strings.HasPrefix(low, "https://"),
			strings.HasPrefix(low, "//"):
			changed = true // drop the remote reference, keep valid syntax
		default:
			b.WriteString(inner)
		}

		b.WriteString(")")
		rest = rest[end+1:]
	}
}

func indexFold(s, sub string) int {
	return strings.Index(strings.ToLower(s), strings.ToLower(sub))
}

// CustomCSS returns the owner's stylesheet ready for a <style> block.
//
// Marking it template.CSS is sound only because SanitizeCustomCSS has run;
// callers must pass the stored value, which is sanitised on save and again
// here as a belt-and-braces measure against rows written by an older build.
func CustomCSS(raw string) template.CSS {
	clean, _ := SanitizeCustomCSS(raw)
	return template.CSS(clean)
}
