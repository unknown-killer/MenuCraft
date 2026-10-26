package server

import (
	"ResturantMenue/internal/models"
	"net/http"
)

// Single-venue deployments retain the editor while exposing a short menu URL.
func (s *Server) menuPath(slug string) string {
	if s.cfg.VenueSlug != "" && slug == s.cfg.VenueSlug {
		return "/menu"
	}
	return "/m/" + slug
}
func (s *Server) handleVenueMenu(w http.ResponseWriter, r *http.Request) {
	if s.cfg.VenueSlug == "" {
		http.NotFound(w, r)
		return
	}
	r.SetPathValue("slug", s.cfg.VenueSlug)
	s.handlePublicMenu(w, r)
}
func (s *Server) handleVenueQR(w http.ResponseWriter, r *http.Request) {
	if s.cfg.VenueSlug == "" {
		http.NotFound(w, r)
		return
	}
	r.SetPathValue("slug", s.cfg.VenueSlug)
	s.handleQRCode(w, r)
}
func (s *Server) handleVenueLanding(w http.ResponseWriter, r *http.Request) {
	venue, err := s.store.Restaurants.BySlug(s.cfg.VenueSlug)
	if err != nil || !venue.IsPublished {
		http.NotFound(w, r)
		return
	}
	v := s.newView(w, r)
	v.Restaurant = venue
	hours, err := s.store.Restaurants.WorkingHours(venue.ID)
	if err != nil {
		s.renderError(w, r, http.StatusInternalServerError)
		return
	}
	v.Page = struct{ Hours []*models.WorkingHour }{Hours: hours}
	s.publicLocale(v, r, venue)
	v.Title = venue.Name(v.Lang())
	s.setVenueSEO(v, false)
	s.render.Render(w, http.StatusOK, "venue", "venue", v)
}
