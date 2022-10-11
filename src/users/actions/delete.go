package useractions

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/mux"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/log"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/session"
	"github.com/abishekmuthian/bonehealthtracker/src/users"
)

func HandleDelete(w http.ResponseWriter, r *http.Request) error {

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

	dexaCookie, err := r.Cookie("reports")

	if err != nil {
		log.Error(log.V{"Upload, Error occurred while reading cookie": err})
	} else {
		log.Info(log.V{"Upload, dexaCookie ": dexaCookie})
	}

	if dexaCookie != nil {
		var dexas []users.Dexa

		report := users.Report{}

		decodedContent, err := base64.StdEncoding.DecodeString(dexaCookie.Value)

		if err == nil {
			err = json.Unmarshal([]byte(decodedContent), &report)

			if err == nil {

				dexas = report.Dexas

				if len(dexas) > 0 {
					dexas = append(dexas[:len(dexas)-1], dexas[len(dexas):]...)

					tempReport := users.Report{
						Dexas: dexas,
					}

					cookieContent, err := json.Marshal(tempReport)

					t, err := time.Parse(time.RFC1123, "Sun, 17 Jan 2038 19:14:07 GMT") // Cookie expires before 2038 bug

					if err != nil {
						log.Error(log.V{"Delete, Setting cookie expires": err})
						return server.InternalError(err)
					}

					if err == nil && len(dexas) > 0 {
						dexaCookie := &http.Cookie{
							Name:    "reports",
							Value:   base64.StdEncoding.EncodeToString(cookieContent),
							Secure:  true,
							Expires: t, // Expires in 2038
							Path:    "/",
						}

						http.SetCookie(w, dexaCookie)
					} else if len(dexas) == 0 {
						// Delete Cookie
						dexaCookie := &http.Cookie{
							Name:    "reports",
							Value:   "",
							Secure:  true,
							Expires: time.Unix(0, 0), // Expires immediately
							Path:    "/",
						}

						http.SetCookie(w, dexaCookie)
					} else {
						log.Error(log.V{"Delete, Error occurred while setting cookie": err})
						return server.InternalError(err)
					}
				} else {
					log.Info(log.V{"Delete": " There are no dexas to delete"})
				}

			}
		} else {
			log.Error(log.V{"Delete, Error occurred while decoding cookie": err})
			return server.InternalError(err)
		}

		// Redirect - ideally here we'd redirect to their original request path
		redirectURL := params.Get("redirectURL")
		if redirectURL == "" {
			redirectURL = "/#upload"
		}

		return server.Redirect(w, r, redirectURL)
	}

	return err
}
