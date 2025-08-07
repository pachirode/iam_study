package docs

import (
	v1 "github.com/pachirode/iam_study/internal/pkg/api/apiserver/v1"
	metaV1 "github.com/pachirode/iam_study/pkg/meta/v1"
)

// swagger:route POST /secrets secrets createSecretRequest
//
// Create secrets.
//
// Create secrets according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: createSecretResponse

// swagger:route DELETE /secrets/{name} secrets deleteSecretRequest
//
// Delete secret.
//
// Delete secret according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: okResponse

// swagger:route PUT /secrets/{name} secrets updateSecretRequest
//
// Update secret.
//
// Update secret according to input parameters.
//
//     Security:
//       api_key:
//
//     Responses:
//       default: errResponse
//       200: updateSecretResponse

// swagger:route GET /secrets/{name} secrets getSecretRequest
//
// Get details for specified secret.
//
// Get details for specified secret according to input parameters.
//
//     Responses:
//       default: errResponse
//       200: getSecretResponse

// swagger:route GET /secrets secrets listSecretRequest
//
// List secrets.
//
// List secrets.
//
//     Responses:
//       default: errResponse
//       200: listSecretResponse

// List users request.
// swagger:parameters listSecretRequest
type listSecretRequestParamsWrapper struct {
	// in:query
	metaV1.ListOptions
}

// List secrets response.
// swagger:response listSecretResponse
type listSecretResponseWrapper struct {
	// in:body
	Body v1.SecretList
}

// Secret response.
// swagger:response createSecretResponse
type createSecretResponseWrapper struct {
	// in:body
	Body v1.Secret
}

// Secret response.
// swagger:response updateSecretResponse
type updateSecretResponseWrapper struct {
	// in:body
	Body v1.Secret
}

// Secret response.
// swagger:response getSecretResponse
type getSecretResponseWrapper struct {
	// in:body
	Body v1.Secret
}

// swagger:parameters createSecretRequest updateSecretRequest
type secretRequestParamsWrapper struct {
	// Secret information.
	// in:body
	Body v1.Secret
}

// swagger:parameters deleteSecretRequest getSecretRequest updateSecretRequest
type secretNameParamsWrapper struct {
	// Secret name.
	// in:path
	Name string `json:"name"`
}
