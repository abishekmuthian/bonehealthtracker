package useractions

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/mux"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/query"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/config"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/log"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/session"
	"github.com/abishekmuthian/bonehealthtracker/src/users"
)

// HandleUpdate handles the update of the research form
// Responds to the post /update
func HandleUpdate(w http.ResponseWriter, r *http.Request) error {

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Get the params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Using turnstile to verify users
	if len(params.Get("cf-turnstile-response")) > 0 {
		if string(params.Get("cf-turnstile-response")) != "" {

			type turnstileResponse struct {
				Success      bool     `json:"success"`
				Challenge_ts string   `json:"challenge_ts"`
				Hostname     string   `json:"hostname"`
				ErrorCodes   []string `json:"error-codes"`
				Action       string   `json:"login"`
				Cdata        string   `json:"cdata"`
			}

			var remoteIP string
			var siteVerify turnstileResponse

			if config.Production() {
				// Get the IP from Cloudflare
				remoteIP = r.Header.Get("CF-Connecting-IP")

			} else {
				// Extract the IP from the address
				remoteIP = r.RemoteAddr
				forward := r.Header.Get("X-Forwarded-For")
				if len(forward) > 0 {
					remoteIP = forward
				}
			}

			postBody := url.Values{}
			postBody.Set("secret", config.Get("turnstile_secret_key"))
			postBody.Set("response", string(params.Get("cf-turnstile-response")))
			postBody.Set("remoteip", remoteIP)

			resp, err := http.Post("https://challenges.cloudflare.com/turnstile/v0/siteverify", "application/x-www-form-urlencoded", strings.NewReader(postBody.Encode()))
			if err != nil {
				log.Info(log.V{"Upload, An error occurred while sending the request to the siteverify": err})
				return server.InternalError(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error(log.V{"Upload, An error occurred while reading the response from the siteverify": err})
				return server.InternalError(err)
			}

			json.Unmarshal(body, &siteVerify)

			if !siteVerify.Success {
				// Security challenge failed
				log.Error(log.V{"Upload, Security challenge failed": siteVerify.ErrorCodes[0]})
				return server.Redirect(w, r, "/?error=security_challenge_failed_research#research")
			}
		} else {
			log.Error(log.V{"Upload, Security challenge unable to process": "response not received from user"})
			return server.Redirect(w, r, "/?error=security_challenge_not_completed_research#research")
		}
	} else {
		// Security challenge not completed
		return server.Redirect(w, r, "/?error=security_challenge_not_completed_research#research")
	}

	if len(params.Get("sex")) > 0 &&
		len(params.Get("first-report-age")) > 0 &&
		len(params.Get("latest-report-age")) > 0 &&
		len(params.Get("treatment")) > 0 &&
		len(params.Get("race-ethnicity")) > 0 {

		var dexas []users.Dexa

		report := users.Report{}

		dexaCookie, err := r.Cookie("reports")

		if err != nil {
			log.Error(log.V{"Update, Error occurred while reading cookie": err})

			return server.InternalError(err, "Submission Failed", "Sorry your submission failed to record contact support")
		} else {
			// Disabling for privacy
			// log.Info(log.V{"Update, dexaCookie ": dexaCookie})
		}

		decodedContent, err := base64.StdEncoding.DecodeString(dexaCookie.Value)

		if err == nil {
			err = json.Unmarshal([]byte(decodedContent), &report)

			if err == nil {

				dexas = report.Dexas

				tempReport := users.Report{
					Sex:             params.Get("sex"),
					FirstReportAge:  int(params.GetInt("first-report-age")),
					LatestReportAge: int(params.GetInt("latest-report-age")),
					Treatment:       params.Get("treatment"),
					RaceEthinicity:  params.Get("race-ethnicity"),
					Dexas:           dexas,
				}

				cookieContent, err := json.Marshal(tempReport)

				if err == nil {

					t, err := time.Parse(time.RFC1123, "Sun, 17 Jan 2038 19:14:07 GMT") // Cookie expires before 2038 bug

					if err != nil {
						log.Error(log.V{"Upload, Setting cookie expires": err})
						return server.InternalError(err)
					}
					dexaCookie := &http.Cookie{
						Name:    "reports",
						Value:   base64.StdEncoding.EncodeToString(cookieContent),
						Secure:  true,
						Expires: t, // Expires in 2038
						Path:    "/",
					}

					http.SetCookie(w, dexaCookie)
				} else {
					log.Error(log.V{"Upload, Error occurred while setting cookie": err})
					return server.InternalError(err)
				}

				// Store in the database
				_, err = query.Exec("insert into reports (report) VALUES($1)", cookieContent)

				if err != nil {
					return server.InternalError(err, "Submission Failed", "Sorry your submission failed to record contact support")
				}

				// Redirect - ideally here we'd redirect to their original request path
				redirectURL := params.Get("redirectURL")
				if redirectURL == "" {
					redirectURL = "/?notice=submission_successful_research#research"
				}

				return server.Redirect(w, r, redirectURL)
			}
		}

	} else {
		// Missing form values
		return server.Redirect(w, r, "/?error=enter_all_required_values_research#research")
	}

	return err
}
