package snap

import "time"

// Standard SNAP HTTP header names.
const (
	HeaderContentType       = "Content-Type"
	HeaderAuthorization     = "Authorization"          // "Bearer <accessToken>"
	HeaderAuthorizationCust = "Authorization-Customer" // B2B2C customer token
	HeaderTimestamp         = "X-TIMESTAMP"
	HeaderSignature         = "X-SIGNATURE"
	HeaderClientKey         = "X-CLIENT-KEY" // access-token request only
	HeaderPartnerID         = "X-PARTNER-ID"
	HeaderExternalID        = "X-EXTERNAL-ID" // numeric, unique per day
	HeaderChannelID         = "CHANNEL-ID"
	HeaderOrigin            = "ORIGIN"
	HeaderIPAddress         = "X-IP-ADDRESS"
	HeaderDeviceID          = "X-DEVICE-ID"
	HeaderLatitude          = "X-LATITUDE"
	HeaderLongitude         = "X-LONGITUDE"
)

// TimestampLayout is the SNAP date/time format: yyyy-MM-ddTHH:mm:ss+HH:mm (RFC3339).
const TimestampLayout = time.RFC3339

// NowTimestamp returns the current time formatted as a SNAP X-TIMESTAMP.
func NowTimestamp() string {
	return time.Now().Format(TimestampLayout)
}
