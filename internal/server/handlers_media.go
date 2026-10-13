package server

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ResturantMenue/internal/models"
)

type mediaPage struct {
	Items []*models.Media
}

func (s *Server) handleMediaPage(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}

	items, err := s.store.Media.List(v.Restaurant.ID, 300)
	if err != nil {
		logf("media: list for venue %d: %v", v.Restaurant.ID, err)
		s.renderError(w, r, http.StatusInternalServerError)
		return
	}

	// The image picker dialog reads the same list as JSON.
	if r.URL.Query().Get("format") == "json" {
		out := make([]map[string]any, 0, len(items))
		for _, m := range items {
			out = append(out, map[string]any{"id": m.ID, "path": m.Path, "name": m.OriginalName})
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": out})
		return
	}

	v.Title = v.T("media.title")
	v.Page = &mediaPage{Items: items}
	s.render.Render(w, http.StatusOK, "panel", "media", v)
}

// imageKind sniffs the leading bytes to identify the format. The browser's
// Content-Type and the filename extension are both attacker-controlled, so
// neither is trusted here.
func imageKind(head []byte) (ext, mime string, ok bool) {
	switch {
	case len(head) >= 3 && bytes.Equal(head[:3], []byte{0xFF, 0xD8, 0xFF}):
		return ".jpg", "image/jpeg", true
	case len(head) >= 8 && bytes.Equal(head[:8], []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}):
		return ".png", "image/png", true
	case len(head) >= 6 && (bytes.Equal(head[:6], []byte("GIF87a")) || bytes.Equal(head[:6], []byte("GIF89a"))):
		return ".gif", "image/gif", true
	case len(head) >= 12 && bytes.Equal(head[:4], []byte("RIFF")) && bytes.Equal(head[8:12], []byte("WEBP")):
		return ".webp", "image/webp", true
	}
	return "", "", false
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "no venue"})
		return
	}

	// Cap the request body before parsing so an oversized upload cannot
	// exhaust memory or disk.
	r.Body = http.MaxBytesReader(w, r.Body, s.cfg.MaxUploadBytes+1<<20)
	if err := r.ParseMultipartForm(4 << 20); err != nil {
		s.uploadFail(w, r, v.T("err.upload_size"), http.StatusRequestEntityTooLarge)
		return
	}
	defer func() {
		if r.MultipartForm != nil {
			_ = r.MultipartForm.RemoveAll()
		}
	}()

	file, header, err := r.FormFile("file")
	if err != nil {
		s.uploadFail(w, r, v.T("err.upload_failed"), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > s.cfg.MaxUploadBytes {
		s.uploadFail(w, r, v.T("err.upload_size"), http.StatusRequestEntityTooLarge)
		return
	}

	head := make([]byte, 512)
	n, _ := io.ReadFull(file, head)
	head = head[:n]

	ext, mime, ok := imageKind(head)
	if !ok {
		s.uploadFail(w, r, v.T("err.upload_type"), http.StatusUnsupportedMediaType)
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		s.uploadFail(w, r, v.T("err.upload_failed"), http.StatusInternalServerError)
		return
	}

	// Files live under a per-venue directory with a random name, so an owner
	// can neither guess nor overwrite another venue's uploads.
	dir := filepath.Join(s.cfg.UploadDir, strconv.FormatInt(v.Restaurant.ID, 10))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logf("upload: mkdir %s: %v", dir, err)
		s.uploadFail(w, r, v.T("err.upload_failed"), http.StatusInternalServerError)
		return
	}

	name, err := randomName(ext)
	if err != nil {
		logf("upload: random name: %v", err)
		s.uploadFail(w, r, v.T("err.upload_failed"), http.StatusInternalServerError)
		return
	}

	dest := filepath.Join(dir, name)
	out, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		logf("upload: create %s: %v", dest, err)
		s.uploadFail(w, r, v.T("err.upload_failed"), http.StatusInternalServerError)
		return
	}

	written, err := io.Copy(out, io.LimitReader(file, s.cfg.MaxUploadBytes))
	closeErr := out.Close()
	if err != nil || closeErr != nil {
		_ = os.Remove(dest)
		logf("upload: write %s: %v / %v", dest, err, closeErr)
		s.uploadFail(w, r, v.T("err.upload_failed"), http.StatusInternalServerError)
		return
	}

	publicPath := "/uploads/" + strconv.FormatInt(v.Restaurant.ID, 10) + "/" + name
	media := &models.Media{
		RestaurantID: v.Restaurant.ID,
		Path:         publicPath,
		OriginalName: truncate(filepath.Base(header.Filename), 120),
		SizeBytes:    written,
		MimeType:     mime,
	}
	if err := s.store.Media.Create(media); err != nil {
		_ = os.Remove(dest)
		logf("upload: record media: %v", err)
		s.uploadFail(w, r, v.T("err.upload_failed"), http.StatusInternalServerError)
		return
	}

	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status": "ok",
			"path":   publicPath,
			"id":     media.ID,
			"name":   media.OriginalName,
		})
		return
	}
	s.setFlash(w, "success", v.T("msg.saved"))
	http.Redirect(w, r, "/panel/media", http.StatusSeeOther)
}

func (s *Server) handleMediaDelete(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}

	id := pathInt(r, "id")
	items, err := s.store.Media.List(v.Restaurant.ID, 1000)
	if err != nil {
		logf("media: list before delete: %v", err)
		s.setFlash(w, "error", v.T("err.generic"))
		http.Redirect(w, r, "/panel/media", http.StatusSeeOther)
		return
	}

	var target *models.Media
	for _, m := range items {
		if m.ID == id {
			target = m
			break
		}
	}
	if target == nil {
		s.setFlash(w, "error", v.T("err.not_found"))
		http.Redirect(w, r, "/panel/media", http.StatusSeeOther)
		return
	}

	if err := s.store.Media.Delete(id, v.Restaurant.ID); err != nil {
		logf("media: delete %d: %v", id, err)
		s.setFlash(w, "error", v.T("err.generic"))
		http.Redirect(w, r, "/panel/media", http.StatusSeeOther)
		return
	}

	// The row is gone; a leftover file is harmless, so a failure here is
	// logged rather than surfaced.
	if disk := s.diskPathFor(v.Restaurant.ID, target.Path); disk != "" {
		if err := os.Remove(disk); err != nil && !os.IsNotExist(err) {
			logf("media: remove file %s: %v", disk, err)
		}
	}

	s.setFlash(w, "success", v.T("msg.deleted"))
	http.Redirect(w, r, "/panel/media", http.StatusSeeOther)
}

// sanitizeUploadPath accepts a stored image path only if it points inside this
// venue's own upload directory. Anything else becomes "", which the templates
// render as "no image".
func (s *Server) sanitizeUploadPath(restaurantID int64, p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}
	prefix := "/uploads/" + strconv.FormatInt(restaurantID, 10) + "/"
	if !strings.HasPrefix(p, prefix) {
		return ""
	}
	name := strings.TrimPrefix(p, prefix)
	if name == "" || strings.ContainsAny(name, `/\`) || strings.Contains(name, "..") {
		return ""
	}
	return p
}

// diskPathFor maps a validated public path back to a file on disk.
func (s *Server) diskPathFor(restaurantID int64, publicPath string) string {
	safe := s.sanitizeUploadPath(restaurantID, publicPath)
	if safe == "" {
		return ""
	}
	return filepath.Join(s.cfg.UploadDir, filepath.FromSlash(strings.TrimPrefix(safe, "/uploads/")))
}

func (s *Server) uploadFail(w http.ResponseWriter, r *http.Request, msg string, status int) {
	if wantsJSON(r) {
		writeJSON(w, status, map[string]string{"error": msg})
		return
	}
	s.setFlash(w, "error", msg)
	http.Redirect(w, r, "/panel/media", http.StatusSeeOther)
}

// randomName produces an unguessable filename with the sniffed extension.
func randomName(ext string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("random: %w", err)
	}
	return hex.EncodeToString(buf) + ext, nil
}
