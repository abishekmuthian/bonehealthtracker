package app

import (
	"net/http"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/config"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/view"
)

// HandlePrivacy displays the home page
// responds to GET /privacy
func HandleLegal(w http.ResponseWriter, r *http.Request) error {

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("meta_title", config.Get("meta_title"))
	view.AddKey("meta_url", config.Get("meta_url"))
	view.AddKey("meta_image", config.Get("meta_image"))
	view.AddKey("meta_title", config.Get("meta_title"))
	view.AddKey("meta_desc", config.Get("meta_desc"))
	view.AddKey("meta_keywords", config.Get("meta_keywords"))
	view.AddKey("meta_twitter", config.Get("meta_twitter"))

	view.Template("app/views/legal.html.got")

	return view.Render()
}
