package useractions

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/mux"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/config"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/log"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/session"
	"github.com/abishekmuthian/bonehealthtracker/src/users"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) error {

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
				return server.Redirect(w, r, "/?error=security_challenge_failed#upload")
			}
		} else {
			log.Error(log.V{"Upload, Security challenge unable to process": "response not received from user"})
			return server.Redirect(w, r, "/?error=security_challenge_not_completed#upload")
		}
	} else {
		// Security challenge not completed
		return server.Redirect(w, r, "/?error=security_challenge_not_completed#upload")
	}

	for _, fh := range params.Files {

		fileType := fh[0].Header.Get("Content-Type")
		fileSize := fh[0].Size

		fileSizeKB := fileSize / 1000

		log.Info(log.V{"Product Submission": "File Upload", "fileType": fileType})
		log.Info(log.V{"Product Submission": "File Upload", "fileSize (kB)": fileSizeKB})

		if fileType == "image/png" || fileType == "image/jpeg" {

			if fileSizeKB > 5000 {
				// Image size is over the limit
				log.Error(log.V{"Upload, Image over size": fileSizeKB})
				return server.Redirect(w, r, "/?error=image_oversize#upload")
			}

			file, err := fh[0].Open()
			defer file.Close()

			if err != nil {
				log.Error(log.V{"Upload, Error opening the uploaded file": err})
			}

			fileData, err := ioutil.ReadAll(file)
			if err != nil {
				log.Error(log.V{"Upload, Error opening the uploaded file": err})

			}

			// Send it to the AWS Comprehend Medical Middleware

			type classifierRequest struct {
				Data []byte `json:"data"`
			}

			type classifierResponse struct {
				Entities []string `json:"Entities"`
			}

			req := &classifierRequest{
				Data: fileData,
			}

			postBody, err := json.Marshal(req)
			if err != nil {
				log.Error(log.V{"Upload, Error creating POST body for classifier": err})
				return server.InternalError(err)
			}

			requestBody := bytes.NewBuffer(postBody)
			resp, err := http.Post(config.Get("classifier_server"), "application/json", requestBody)
			if err != nil {
				log.Info(log.V{"Upload, An error occurred while sending the request to the middleware": err})
				return server.InternalError(err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error(log.V{"Upload, An error occurred while reading the response from the middleware": err})
				return server.InternalError(err)
			}

			organs := users.Parse(body)

			if len(organs) > 0 {
				// Disabling log for privacy
				// log.Info(log.V{"Organs after parsing ": organs})

				dexaCookie, err := r.Cookie("reports")

				if err != nil {
					log.Error(log.V{"Upload, Cookie not found error, Probably no report was uploaded yet": err})
				} else {
					// Disabling log for privacy
					// log.Info(log.V{"Upload, dexaCookie ": dexaCookie})
				}

				if dexaCookie == nil {

					var dexas []users.Dexa

					dexa := users.Dexa{
						Id:     1,
						Year:   "Year 1",
						Organs: organs,
					}

					dexas = append(dexas, dexa)

					report := users.Report{
						Dexas: dexas,
					}

					cookieContent, err := json.Marshal(report)

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

				} else {
					var dexas []users.Dexa

					report := users.Report{}

					decodedContent, err := base64.StdEncoding.DecodeString(dexaCookie.Value)

					if err == nil {
						err = json.Unmarshal([]byte(decodedContent), &report)

						if err == nil {

							id := len(report.Dexas) + 1

							tempDexa := users.Dexa{
								Id:     id,
								Year:   "Year " + strconv.Itoa(id),
								Organs: organs,
							}

							dexas = report.Dexas

							dexas = append(dexas, tempDexa)

							tempReport := users.Report{
								Dexas: dexas,
							}

							cookieContent, err := json.Marshal(tempReport)

							// Check if the cookie size limit has been reached, Accounting for the default size
							if cookieContentSize := len(base64.StdEncoding.EncodeToString(cookieContent)); cookieContentSize > 4096 {
								log.Info(log.V{"Upload, Cookie size limit reached": cookieContentSize})

								// Cookie size limit reached
								return server.Redirect(w, r, "/?error=max_reports_reached#upload")
							}

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
						}
					} else {
						log.Error(log.V{"Upload, Error occurred while decoding cookie": err})
						return server.InternalError(err)
					}

				}
			} else {
				// No organs found
				return server.Redirect(w, r, "/?error=not_a_valid_report#upload")
			}

			// Redirect - ideally here we'd redirect to their original request path
			redirectURL := params.Get("redirectURL")
			if redirectURL == "" {
				redirectURL = "/#reports"
			}

			return server.Redirect(w, r, redirectURL)

		} else {
			// Not a valid image
			return server.Redirect(w, r, "/?error=not_a_valid_report#upload")
		}

	}

	return err

}
