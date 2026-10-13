package theme

import "strings"

import "testing"

// TestSanitizeCustomCSSBlocksEscapes covers the ways a stylesheet could stop
// being a stylesheet. Each of these would otherwise reach a <style> block.
func TestSanitizeCustomCSSBlocksEscapes(t *testing.T) {
	cases := []struct {
		name string
		in   string
		gone string
	}{
		{"closing style tag", `.a{}</style><script>alert(1)</script>`, "</style"},
		{"uppercase closing tag", `.a{}</STYLE><img onerror=x>`, "</style"},
		{"html comment", `.a{} <!-- .b{} -->`, "<!--"},
		{"import", `@import url("https://evil.example/x.css"); .a{}`, "@import"},
		{"expression", `.a{width:expression(alert(1))}`, "expression("},
		{"javascript url", `.a{background:url(javascript:alert(1))}`, "javascript:"},
		{"moz binding", `.a{-moz-binding:url(evil.xml)}`, "-moz-binding"},
		{"script tag", `.a{}<script>x</script>`, "<script"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, removed := SanitizeCustomCSS(tc.in)

			if strings.Contains(strings.ToLower(out), tc.gone) {
				t.Errorf("%q survived sanitising: %q", tc.gone, out)
			}
			if len(removed) == 0 {
				t.Errorf("sanitiser did not report removing anything from %q", tc.in)
			}
			// The whole line goes, so no mangled fragment is left behind.
			for _, fragment := range []string{"alert(", "onerror", "evil.example"} {
				if strings.Contains(strings.ToLower(out), fragment) {
					t.Errorf("fragment %q left behind after removing the line: %q", fragment, out)
				}
			}
		})
	}
}

// TestSanitizeCustomCSSDropsRemoteFetches checks the privacy rule: a stylesheet
// must not be able to pull assets from a third party, which would leak every
// visitor's IP to that host.
func TestSanitizeCustomCSSDropsRemoteFetches(t *testing.T) {
	in := `.hero{background:url("https://tracker.example/pixel.png")}
.logo{background:url('//cdn.example/x.png')}
.ok{background:url("/uploads/2/photo.jpg")}
.fine{background:url(data:image/gif;base64,R0lGOD)}`

	out, removed := SanitizeCustomCSS(in)

	if strings.Contains(out, "tracker.example") || strings.Contains(out, "cdn.example") {
		t.Errorf("a remote url() survived: %q", out)
	}
	if !strings.Contains(out, "/uploads/2/photo.jpg") {
		t.Errorf("a same-origin upload url() was wrongly dropped: %q", out)
	}
	if !strings.Contains(out, "data:image/gif") {
		t.Errorf("a data: url() was wrongly dropped: %q", out)
	}
	if len(removed) == 0 {
		t.Error("removing remote urls was not reported")
	}
}

// TestSanitizeCustomCSSKeepsRealStyling makes sure the filter is not so blunt
// that the feature stops being useful.
func TestSanitizeCustomCSSKeepsRealStyling(t *testing.T) {
	in := `.m-card {
  border-radius: 28px;
  box-shadow: 0 10px 40px rgba(0,0,0,.3);
  transition: transform .3s ease;
}
.m-card:hover { transform: translateY(-6px) scale(1.02); }
@media (max-width: 600px) { .m-grid { grid-template-columns: 1fr; } }
@keyframes float { 50% { transform: translateY(-8px); } }
.m-name { font-family: "My Font", serif; letter-spacing: .04em; }`

	out, removed := SanitizeCustomCSS(in)

	if len(removed) != 0 {
		t.Errorf("legitimate CSS was filtered: %v", removed)
	}
	for _, want := range []string{"border-radius", "@media", "@keyframes", ":hover", "translateY"} {
		if !strings.Contains(out, want) {
			t.Errorf("%q was lost from a normal stylesheet", want)
		}
	}
}

// TestSanitizeCustomCSSCapsSize stops one venue's stylesheet from bloating
// every page render.
func TestSanitizeCustomCSSCapsSize(t *testing.T) {
	huge := strings.Repeat(".a{color:red}", 5_000)
	out, removed := SanitizeCustomCSS(huge)

	if len(out) > MaxCustomCSS {
		t.Errorf("output is %d bytes, over the %d cap", len(out), MaxCustomCSS)
	}
	if len(removed) == 0 {
		t.Error("truncation was not reported")
	}
}

// TestSanitizeCustomCSSEmpty is the common case and must stay cheap and quiet.
func TestSanitizeCustomCSSEmpty(t *testing.T) {
	if out, removed := SanitizeCustomCSS("   "); out != "" || removed != nil {
		t.Errorf("empty input produced %q / %v", out, removed)
	}
}
