package server

import (
	"bytes"
	"encoding/base64"
	"html"
	"image/png"
	"net/http"
	"regexp"
	"strings"
	"testing"

	qrcode "github.com/skip2/go-qrcode"
)

func TestTableQRLink(t *testing.T) {
	s, _, v := fixture(t)
	w := request(s, "GET", "/m/"+v.Slug+"/qr?table=۰۱۲&size=512&download=1", nil, nil)
	if w.Code != 200 {
		t.Fatalf("status %d", w.Code)
	}
	want, err := qrcode.Encode("https://menu.example.com/m/"+v.Slug+"?table=12", qrcode.Medium, 512)
	if err != nil || !bytes.Equal(w.Body.Bytes(), want) {
		t.Fatal("QR does not contain the normalized table link", err)
	}
	if !strings.Contains(w.Header().Get("Content-Disposition"), "-table-12.png") {
		t.Fatal("download filename loses table")
	}
	for _, value := range []string{"0", "-1", "10000", "1.2", "1e2", "abc"} {
		if got := request(s, "GET", "/m/"+v.Slug+"/qr?table="+value, nil, nil); got.Code != 400 {
			t.Fatalf("accepted table %s", value)
		}
	}
	page := request(s, "GET", "/m/"+v.Slug+"?table=۱۲", nil, nil).Body.String()
	if !strings.Contains(page, "table=12") || !strings.Contains(page, "data-notebook=") || !strings.Contains(page, "data-note-item=") {
		t.Fatal("table or notebook missing from public menu")
	}
}

func TestStandLayoutsAndOwnerScope(t *testing.T) {
	s, u, v := fixture(t)
	if w := request(s, "GET", "/panel/qr", nil, nil); w.Code != http.StatusSeeOther {
		t.Fatal("stand editor must require login")
	}
	for _, lang := range []string{"fa", "en"} {
		for _, layout := range []string{"above", "below", "left", "right", "tent"} {
			w := request(s, "GET", "/panel/qr?start=۱۲&end=۱۳&layout="+layout+"&lang="+lang, nil, session(s, u))
			body := html.UnescapeString(w.Body.String())
			if w.Code != 200 || strings.Contains(body, "ZgotmplZ") {
				t.Fatalf("%s/%s: %d %s", lang, layout, w.Code, body[:min(len(body), 400)])
			}
			if !strings.Contains(body, v.Name(lang)) || !strings.Contains(body, "table=12") || !strings.Contains(body, "table=13") {
				t.Fatal("wrong venue or table links")
			}
			images := regexp.MustCompile(`src="data:image/png;base64,([A-Za-z0-9+/=]+)"`).FindAllStringSubmatch(body, -1)
			count := 2
			if layout == "tent" {
				count = 4
			}
			if len(images) != count {
				t.Fatalf("%s: %d QR images", layout, len(images))
			}
			data, err := base64.StdEncoding.DecodeString(images[0][1])
			if err != nil {
				t.Fatal(err)
			}
			img, err := png.Decode(bytes.NewReader(data))
			if err != nil || img.Bounds().Dx() != 768 {
				t.Fatal("invalid print image", err)
			}
		}
	}
	for _, query := range []string{"start=1&end=101", "start=8&end=2", "start=0", "end=10000", "layout=invalid"} {
		if w := request(s, "GET", "/panel/qr?"+query, nil, session(s, u)); w.Code != 400 {
			t.Fatal("invalid batch accepted", query)
		}
	}
}
