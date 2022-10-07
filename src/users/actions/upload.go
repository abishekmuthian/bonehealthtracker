package useractions

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

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
		return err
	}

	// Check if this IP is allowed to upload

	// Get the params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	for _, fh := range params.Files {

		fileType := fh[0].Header.Get("Content-Type")
		fileSize := fh[0].Size

		log.Info(log.V{"Product Submission": "File Upload", "fileType": fileType})
		log.Info(log.V{"Product Submission": "File Upload", "fileSize (kB)": fileSize / 1000})

		if fileType == "image/png" || fileType == "image/jpeg" {
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
			}

			responseBody := bytes.NewBuffer(postBody)
			resp, err := http.Post(config.Get("classifier_server"), "application/json", responseBody)
			if err != nil {
				log.Info(log.V{"Upload, An error occurred while sending the request to the middleware": err})
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Error(log.V{"Upload, An error occurred while reading the response from the middleware": err})
			}

			organs := users.Parse(body)

			if organs != nil {
				log.Info(log.V{"Organs after parsing ": organs})

				dexaCookie, err := r.Cookie("reports")

				if err != nil {
					log.Error(log.V{"Upload, Error occurred while reading cookie": err})
				} else {
					log.Info(log.V{"Upload, dexaCookie ": dexaCookie})
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
						dexaCookie := &http.Cookie{
							Name:   "reports",
							Value:  base64.StdEncoding.EncodeToString(cookieContent),
							MaxAge: 86400 * 60,
							Path:   "/",
						}

						http.SetCookie(w, dexaCookie)
					} else {
						log.Error(log.V{"Upload, Error occurred while setting cookie": err})
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

							if err == nil {
								dexaCookie := &http.Cookie{
									Name:   "reports",
									Value:  base64.StdEncoding.EncodeToString(cookieContent),
									MaxAge: 86400 * 60,
									Path:   "/",
								}

								http.SetCookie(w, dexaCookie)
							} else {
								log.Error(log.V{"Upload, Error occurred while setting cookie": err})
							}
						}
					} else {
						log.Error(log.V{"Upload, Error occurred while decoding cookie": err})
					}

				}
			}

			// Redirect - ideally here we'd redirect to their original request path
			redirectURL := params.Get("redirectURL")
			if redirectURL == "" {
				redirectURL = "/"
			}

			return server.Redirect(w, r, redirectURL)

		} else {

		}

	}

	return err

}
