/* cors enables CORS for all API endpoints that use it. */

package cors

import (
	"net/http"
)

// CORS simply enables CORS for an endpoint
func CORS(wRef *http.ResponseWriter) {
	(*wRef).Header().Set("Content-Type", "application/json")
	(*wRef).Header().Set("Access-Control-Allow-Origin", "*")
	(*wRef).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*wRef).Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Access-Control-Allow-Methods, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Token")
}
