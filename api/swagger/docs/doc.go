// Package docs IAM API Server API.
//
// Identity and Access Management System.
//
//	    Schemes: http, https
//	    BasePath: /v1
//	    Version: 1.0.0
//		   License: MIT https://opensource.org/licenses/MIT
//
//	    Consumes:
//	    - application/json
//
//	    Produces:
//	    - application/json
//
//	    Security:
//	    - basic
//	    - api_key
//
//	   SecurityDefinitions:
//	   basic:
//	     type: basic
//	   api_key:
//	     type: apiKey
//	     name: Authorization
//	     in: header
//
// swagger:meta
package docs

import "github.com/pachirode/iam_study/pkg/core"

// ErrResponse defines the return messages when an error occurred.
// swagger:response errResponse
type errResponseWrapper struct {
	// in:body
	Body core.ErrResponse
}

// Return nil json object.
// swagger:response okResponse
type okResponseWrapper struct{}
