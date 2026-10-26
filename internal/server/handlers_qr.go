package server

import (
	"encoding/base64"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// A table label is display context, not an authenticated order or reservation.
func validTableNumber(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.Map(func(r rune) rune {
		if r >= '۰' && r <= '۹' {
			return '0' + r - '۰'
		}
		if r >= '٠' && r <= '٩' {
			return '0' + r - '٠'
		}
		return r
	}, raw)
	if len(raw) == 0 || len(raw) > 4 {
		return ""
	}
	for _, r := range raw {
		if r < '0' || r > '9' {
			return ""
		}
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 || n > 9999 {
		return ""
	}
	return strconv.Itoa(n)
}

func tableMenuURL(base, table string) string {
	if table == "" {
		return base
	}
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	q := u.Query()
	q.Set("table", table)
	u.RawQuery = q.Encode()
	return u.String()
}

type qrStand struct {
	Number int
	Image  template.URL
	URL    string
}
type qrStandPage struct {
	Start, End, Layout string
	Sheets             [][]qrStand
	Error              string
	Published          bool
}

func (s *Server) handleQRStands(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		s.renderError(w, r, http.StatusNotFound)
		return
	}
	p := &qrStandPage{Start: "1", End: "1", Layout: "above", Published: v.Restaurant.IsPublished}
	q := r.URL.Query()
	if q.Has("start") {
		p.Start = validTableNumber(q.Get("start"))
	}
	if q.Has("end") {
		p.End = validTableNumber(q.Get("end"))
	}
	if q.Has("layout") {
		p.Layout = q.Get("layout")
	}
	start, _ := strconv.Atoi(p.Start)
	end, _ := strconv.Atoi(p.End)
	if start == 0 || end == 0 || end < start || end-start >= 100 {
		p.Error = v.T("qr.range_error")
	}
	switch p.Layout {
	case "above", "below", "left", "right", "tent":
	default:
		p.Error = v.T("qr.layout_error")
		p.Layout = "above"
	}
	if p.Error == "" {
		perSheet := 4
		if p.Layout == "tent" {
			perSheet = 2
		}
		for n := start; n <= end; n++ {
			number := strconv.Itoa(n)
			link := tableMenuURL(s.absoluteMenuURL(r, v.Restaurant.Slug), number)
			png, err := qrcode.Encode(link, qrcode.Medium, 768)
			if err != nil {
				s.renderError(w, r, http.StatusInternalServerError)
				return
			}
			if (n-start)%perSheet == 0 {
				p.Sheets = append(p.Sheets, []qrStand{})
			}
			last := len(p.Sheets) - 1
			p.Sheets[last] = append(p.Sheets[last], qrStand{Number: n, Image: template.URL("data:image/png;base64," + base64.StdEncoding.EncodeToString(png)), URL: link})
		}
	}
	v.Title = v.T("qr.title")
	v.Page = p
	w.Header().Set("Cache-Control", "private, no-store")
	status := http.StatusOK
	if p.Error != "" {
		status = http.StatusBadRequest
	}
	s.render.Render(w, status, "qr", "qr", v)
}
