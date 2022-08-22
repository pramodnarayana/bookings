package handlers

import (
	"net/http"

	"github.com/pramodnarayana/bookings/pkg/config"
	"github.com/pramodnarayana/bookings/pkg/models"
	"github.com/pramodnarayana/bookings/pkg/render"
)

// Repo the repository used by the handlers
var Repo Repository

// Repository is the respository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = *r
}

// This is Home Page
func (repository *Repository) Home(rw http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	repository.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(rw, "home.page.html", &models.TemplateData{})

}

// This is About Page
func (repository *Repository) About(rw http.ResponseWriter, r *http.Request) {
	// Perform some logic
	stringMap := map[string]string{}
	stringMap["test"] = "Hello again"

	remoteIP := repository.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// Send the data to the template
	render.RenderTemplate(rw, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
