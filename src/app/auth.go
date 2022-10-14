package app

import (
	"github.com/abishekmuthian/bonehealthtracker/src/lib/auth"
	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/config"
)

// SetupAuth sets up the auth pkg and authorisation for users
func SetupAuth() {

	// Set up the auth package with our secrets from config
	auth.HMACKey = auth.HexToBytes(config.Get("hmac_key"))
	auth.SecretKey = auth.HexToBytes(config.Get("secret_key"))
	auth.SessionName = config.Get("session_name")

	// Enable https cookies on production server - everyone should be on https
	if config.Production() {
		auth.SecureCookies = true
	}

	// Set up our authorisation for user roles on resources using can pkg

	// No authorizations required currently for Bone Health Tracker

}
