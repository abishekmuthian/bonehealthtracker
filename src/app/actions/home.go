package productctions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/mux"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/config"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/log"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/view"
	"github.com/abishekmuthian/bonehealthtracker/src/users"
)

// HandleHome displays the home page
// responds to GET /
func HandleHome(w http.ResponseWriter, r *http.Request) error {

	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)

	view.AddKey("meta_title", config.Get("meta_title"))
	view.AddKey("meta_url", config.Get("meta_url"))
	view.AddKey("meta_image", config.Get("meta_image"))
	view.AddKey("meta_title", config.Get("meta_title"))
	view.AddKey("meta_desc", config.Get("meta_desc"))
	view.AddKey("meta_keywords", config.Get("meta_keywords"))
	view.AddKey("meta_twitter", config.Get("meta_twitter"))

	view.Template("app/views/home.html.got")

	view.AddKey("error", params.Get("error"))
	view.AddKey("notice", params.Get("notice"))

	// Set Cloudflare turnstile site key
	view.AddKey("turnstile_site_key", config.Get("turnstile_site_key"))

	dexaCookie, err := r.Cookie("reports")

	if err != nil {
		log.Error(log.V{"Home, Cookie not found error, Probably no report was uploaded yet": err})
	} else {
		// Disabling log for privacy
		// log.Info(log.V{"Home, dexaCookie ": dexaCookie})
	}

	if dexaCookie != nil {

		report := users.Report{}

		decodedContent, err := base64.StdEncoding.DecodeString(dexaCookie.Value)

		if err == nil {
			err = json.Unmarshal([]byte(decodedContent), &report)

			if err == nil {
				view.AddKey("report", report)

				skeletonMap := make(map[string]float64)

				dexa := report.Dexas[len(report.Dexas)-1]

				for _, organ := range dexa.Organs {

					site := strings.ToLower(organ.Site)
					direction := strings.ToLower(organ.Direction)

					if strings.Contains(site, "spine") || strings.Contains(site, "l1-l4") || strings.Contains(site, "l1 through l4") {
						skeletonMap["apSpine"] = organ.TScore
					}

					if site == "femur" {
						if direction == "left" {
							skeletonMap["leftFemur"] = organ.TScore
						} else {
							skeletonMap["rightFemur"] = organ.TScore
						}
					}

					if site == "femur neck" || site == "femoral neck" {
						if direction == "left" {
							skeletonMap["leftFemurNeck"] = organ.TScore
						} else {
							skeletonMap["rightFemurNeck"] = organ.TScore
						}
					}

					if site == "hip" {
						if direction == "left" {
							skeletonMap["leftHip"] = organ.TScore
						} else {
							skeletonMap["rightHip"] = organ.TScore
						}
					}

					if site == "forearm" {
						skeletonMap["forearm"] = organ.TScore
					}

				}

				view.AddKey("skeletonMap", skeletonMap)

			} else {
				log.Error(log.V{"Home, Error while marshalling cookie ": err})
			}

		}
	}

	return view.Render()
}
