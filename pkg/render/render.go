package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/pramodnarayana/bookings/pkg/config"
	"github.com/pramodnarayana/bookings/pkg/models"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(defaultTemplateData *models.TemplateData) *models.TemplateData {
	return defaultTemplateData
}

// Render templates dynamically
func RenderTemplate(rw http.ResponseWriter, templateName string, templateData *models.TemplateData) {

	var templateCache map[string]*template.Template
	if app.UseCache {
		// Get the template cache from the AppConfig
		templateCache = app.TemplateCache
	} else {
		// Create template cache
		templateCache, _ = CreateParsedTemplateCache()
	}

	// Get requested template from the cache
	parsedTemplate, ok := templateCache[templateName]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	err := parsedTemplate.Execute(buf, templateData)
	if err != nil {
		log.Println(err)
	}

	// Render the template
	_, err = buf.WriteTo(rw)
	if err != nil {
		log.Println("Error writing template to the browser", err)
	}
}

func CreateParsedTemplateCache() (map[string]*template.Template, error) {
	myCache := make(map[string]*template.Template)

	// Get all of the files name *.page.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	const layoutPattern = "./templates/*.layout.html"
	// Range through all the files ending with *.page.html
	for _, page := range pages {
		name := filepath.Base(page)
		parsedTemplate, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		layoutTemplates, err := filepath.Glob(layoutPattern)
		if err != nil {
			return myCache, err
		}

		if len(layoutTemplates) > 0 {
			parsedTemplate, err = parsedTemplate.ParseGlob(layoutPattern)
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = parsedTemplate
	}

	return myCache, nil
}

/* var parsedTemplateCache = make(map[string]*template.Template)

func RenderTemplate(rw http.ResponseWriter, templateName string) {
	var parsedTemplate *template.Template
	var err error

	//check to see if we already have templates  in our cache
	_, inMap := parsedTemplateCache[templateName]
	if !inMap {
		// Need to create template
		log.Println("Creating template and adding to cache")
		err = createParsedTemplateCache(templateName)
		if err != nil {
			log.Println(err)
		}
	} else {
		// We have the template in the cache
		log.Println("Using cached template")
		parsedTemplate = parsedTemplateCache[templateName]
	}

	err = parsedTemplate.Execute(rw, nil)
	if err != nil {
		log.Println(err)
	}

}

func createParsedTemplateCache(templateName string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", templateName),
		"./templates/base.layout.html",
	}

	// Parse template
	parsedTemplate, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// Add parsed template to the cache
	parsedTemplateCache[templateName] = parsedTemplate
	return nil

} */
