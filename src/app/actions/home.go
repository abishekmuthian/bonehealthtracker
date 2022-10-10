package productctions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/auth"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/mux"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/config"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/log"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/session"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/view"
	"github.com/abishekmuthian/bonehealthtracker/src/users"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/stats"
)

// HandleHome displays the home page
// responds to GET /
func HandleHome(w http.ResponseWriter, r *http.Request) error {
	stats.RegisterHit(r)

	currentUser := session.CurrentUser(w, r)

	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("home", 1)

	view.AddKey("meta_title", config.Get("meta_title"))
	view.AddKey("meta_url", config.Get("meta_url"))
	view.AddKey("meta_image", config.Get("meta_image"))
	view.AddKey("meta_title", config.Get("meta_title"))
	view.AddKey("meta_desc", config.Get("meta_desc"))
	view.AddKey("meta_keywords", config.Get("meta_keywords"))
	view.AddKey("meta_twitter", config.Get("meta_twitter"))

	view.Template("app/views/home.html.got")

	view.AddKey("currentUser", currentUser)

	view.AddKey("error", params.Get("error"))
	view.AddKey("notice", params.Get("notice"))
	view.AddKey("show_reset_password", params.Get("show_reset_password"))
	view.AddKey("email", params.Get("email"))

	view.AddKey("userCount", stats.UserCount())

	view.AddKey("itemId", params.Get("itemId"))

	if currentUser.Anon() {
		view.AddKey("loggedIn", false)
	} else {
		view.AddKey("loggedIn", true)
		view.AddKey("clientID", config.Get("client_Id"))
		view.AddKey("redirectURI", config.Get("twitter_redirect_uri"))
		view.AddKey("twitterScopes", config.Get("twitter_scopes"))

		nonceToken, err := auth.NonceToken(w, r)

		if err == nil {
			view.AddKey("code", nonceToken)
		}
	}

	//view.AddKey("validationDeadline", math.Round(time.Date(2021, time.June, 30, 0, 0, 0, 0, time.UTC).Sub(time.Now()).Hours()/24))

	if currentUser.Anon() {
		view.AddKey("disableSubmit", false)
		view.AddKey("disableCopy", true)
	} else {
		view.AddKey("disableSubmit", true)
		view.AddKey("disableCopy", false)

		clientCountry := r.Header.Get("CF-IPCountry")
		log.Info(log.V{"Subscription, Client Country": clientCountry})
		if !config.Production() {
			// There will be no CF request header in the development/test
			clientCountry = config.Get("subscription_client_country")
		}

		if clientCountry == "IN" {
			view.AddKey("priceId", config.Get("stripe_price_id_ideator_IN"))
			view.AddKey("price", config.Get("stripe_price_IN"))
		} else if clientCountry == "GB" {
			view.AddKey("priceId", config.Get("stripe_price_id_ideator_GB"))
			view.AddKey("price", config.Get("stripe_price_GB"))
		} else if clientCountry == "CA" {
			view.AddKey("priceId", config.Get("stripe_price_id_ideator_CA"))
			view.AddKey("price", config.Get("stripe_price_CA"))
		} else if clientCountry == "AU" {
			view.AddKey("priceId", config.Get("stripe_price_id_ideator_AU"))
			view.AddKey("price", config.Get("stripe_price_AU"))
		} else if clientCountry == "DE" || clientCountry == "FR" {
			view.AddKey("priceId", config.Get("stripe_price_id_ideator_EU"))
			view.AddKey("price", config.Get("stripe_price_EU"))
		} else {
			view.AddKey("priceId", config.Get("stripe_price_id_ideator_US"))
			view.AddKey("price", config.Get("stripe_price_US"))
		}
	}

	dexaCookie, err := r.Cookie("reports")

	if err != nil {
		log.Error(log.V{"Home, Error occurred while reading cookie": err})
	} else {
		log.Info(log.V{"Home, dexaCookie ": dexaCookie})
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

				}

				view.AddKey("skeletonMap", skeletonMap)

			} else {
				log.Error(log.V{"Home, Error while marshalling cookie ": err})
			}

		}
	}

	return view.Render()
}
