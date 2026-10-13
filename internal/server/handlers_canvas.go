package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"ResturantMenue/internal/models"
	"ResturantMenue/internal/repository"
)

type canvasPage struct {
	Block      *models.Block
	Canvas     models.CanvasData
	CanvasJSON string
	Kinds      []models.CanvasElementKind
	IsNew      bool
}

// handleCanvasEdit opens the free-form designer for one canvas block.
func (s *Server) handleCanvasEdit(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}

	id := pathInt(r, "id")
	page := &canvasPage{Kinds: models.CanvasElementKinds}

	if id > 0 {
		b, err := s.store.Blocks.ByID(id, v.Restaurant.ID)
		if err != nil {
			if !errors.Is(err, repository.ErrNotFound) {
				logf("canvas: load block %d: %v", id, err)
			}
			s.renderError(w, r, http.StatusNotFound)
			return
		}
		page.Block = b
		page.Canvas = b.Canvas
	} else {
		page.IsNew = true
		page.Canvas = models.DefaultCanvas()
	}

	page.CanvasJSON = page.Canvas.CanvasJSON()

	v.Title = v.T("cv.title")
	v.Page = page
	s.render.Render(w, http.StatusOK, "panel", "canvas", v)
}

// handleCanvasSave persists the designed canvas.
func (s *Server) handleCanvasSave(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Redirect(w, r, "/panel/", http.StatusSeeOther)
		return
	}
	if err := parseForm(r); err != nil {
		s.blockRedirect(w, r, "error", v.T("err.generic"))
		return
	}

	data := models.ParseCanvas(formStr(r, "canvas"))
	// Background images must belong to this venue.
	data.BgImage = s.sanitizeUploadPath(v.Restaurant.ID, data.BgImage)
	for i := range data.Elements {
		data.Elements[i].Image = s.sanitizeUploadPath(v.Restaurant.ID, data.Elements[i].Image)
	}

	id, _ := strconv.ParseInt(formStr(r, "id"), 10, 64)

	if id > 0 {
		b, err := s.store.Blocks.ByID(id, v.Restaurant.ID)
		if err != nil {
			s.blockRedirect(w, r, "error", v.T("err.not_found"))
			return
		}
		b.Canvas = data
		b.TitleFa = formStr(r, "title_fa")
		b.TitleEn = formStr(r, "title_en")
		b.IsVisible = formBool(r, "is_visible")
		if err := s.store.Blocks.Update(b); err != nil {
			logf("canvas: update block %d: %v", id, err)
			s.blockRedirect(w, r, "error", v.T("err.generic"))
			return
		}
		s.blockRedirect(w, r, "success", v.T("msg.saved"))
		return
	}

	b := &models.Block{
		RestaurantID: v.Restaurant.ID,
		Kind:         models.BlockCanvas,
		TitleFa:      formStr(r, "title_fa"),
		TitleEn:      formStr(r, "title_en"),
		IsVisible:    true,
		Canvas:       data,
		Style:        models.DefaultBlockStyle(models.BlockCanvas),
	}
	// A canvas paints its own background and spans the page.
	b.Style.Width = "full"
	b.Style.Padding = "none"

	if err := s.store.Blocks.Create(b); err != nil {
		logf("canvas: create for venue %d: %v", v.Restaurant.ID, err)
		s.blockRedirect(w, r, "error", v.T("err.generic"))
		return
	}

	s.appendBlockToSections(v.Restaurant.ID, b.ID)
	s.blockRedirect(w, r, "success", v.T("msg.created"))
}

// handleCanvasRender returns the canvas markup for the editor stage.
//
// The editor calls this after any styling change so the surface being edited is
// produced by exactly the same Go code that renders the public page.
func (s *Server) handleCanvasRender(w http.ResponseWriter, r *http.Request) {
	v := s.newView(w, r)
	if v.Restaurant == nil {
		http.Error(w, "no venue", http.StatusBadRequest)
		return
	}
	if err := parseForm(r); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	data := models.ParseCanvas(formStr(r, "canvas"))
	data.BgImage = s.sanitizeUploadPath(v.Restaurant.ID, data.BgImage)
	for i := range data.Elements {
		data.Elements[i].Image = s.sanitizeUploadPath(v.Restaurant.ID, data.Elements[i].Image)
	}

	t, err := s.render.page("panel", "canvas")
	if err != nil {
		httpError(w, err)
		return
	}

	var buf bytes.Buffer
	payload := map[string]any{
		"cv": CanvasView{Data: data, Locale: v.Locale.Code, Edit: true},
		"v":  v,
	}
	if err := t.ExecuteTemplate(&buf, "canvas", payload); err != nil {
		httpError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = buf.WriteTo(w)
}

// handleCanvasElement returns a new element of the requested kind, with the
// defaults defined in Go rather than duplicated in the editor.
func (s *Server) handleCanvasElement(w http.ResponseWriter, r *http.Request) {
	if err := parseForm(r); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	kind := formStr(r, "kind")
	if !models.ValidCanvasKind(kind) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown kind"})
		return
	}

	id := formStr(r, "id")
	if id == "" {
		id = "el"
	}

	e := models.DefaultCanvasElement(models.CanvasElementKind(kind), id)
	e.Normalize()

	out, err := json.Marshal(e)
	if err != nil {
		httpError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(out)
}
